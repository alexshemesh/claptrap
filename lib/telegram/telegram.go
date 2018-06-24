package telegram

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"

	"github.com/alexshemesh/claptrap/lib/types"
	"encoding/json"
)

var settingsKey = "telegram/settings"

type TelegramBot struct{
	log             logs.Logger
	token           string
	bot             *tgbotapi.BotAPI
	settingsHandler types.Settings
}


func NewTelegramBot(logPar logs.Logger, settingsHandlerPar types.Settings)(retVal *TelegramBot) {
	tokenPar,err := settingsHandlerPar.GetValue("telegram/token" )
	if err == nil {
		retVal = &TelegramBot{log: logPar, token: tokenPar, settingsHandler: settingsHandlerPar }
	}
	return retVal
}

func (this *TelegramBot)Connect()(err error){
	this.bot, err = tgbotapi.NewBotAPI(this.token)
	return err
}

func (this TelegramBot)getSettings()(retVal types.TelegramSettings, err error){
	var data string
	data,err = this.settingsHandler.GetValue(settingsKey)
	if err == nil && data  != "" {
		err = json.Unmarshal([]byte(data), &retVal)
	}
	return retVal,err
}

func (this TelegramBot)GetMessages(  )(retVal []tgbotapi.Message, err error) {
	var settings types.TelegramSettings
	settings, err = this.getSettings( )
	if err == nil || err.Error() == "404 Not Found"{
		u := tgbotapi.NewUpdate(settings.LastUpdate + 1)
		var updates []tgbotapi.Update
		if this.bot != nil {
			updates, err = this.bot.GetUpdates(u)

			if err == nil {
				for _, update := range updates {
					settings.LastUpdate = update.UpdateID
					this.log.Debug(fmt.Sprintf("Update id %d", update.UpdateID))
					if update.Message != nil {
						retVal = append(retVal, *update.Message)
					}
				}
			}
		} else {
			err = fmt.Errorf("Not connected")
			this.log.Error(err)
		}
		dataToSend,_ := json.Marshal(settings)
		this.settingsHandler.SetValue(settingsKey, string( dataToSend ) )
	}
	return retVal, err
}

func (this TelegramBot) Ping()(err error){
	if this.bot.Self.UserName == "" {
		err = fmt.Errorf("Bot is not connected")
	}
	this.log.Debug("Telegram Ping is OK. Connected with username " + this.bot.Self.UserName)
	return err
}

func GetMsgContactTitle(msg tgbotapi.Message )(retVal string){
	retVal = "Unknown"
	if msg.Contact != nil {
		retVal = fmt.Sprintf("%s %s",msg.Contact.FirstName , msg.Contact.LastName)
	}else if msg.Chat != nil {
		retVal = fmt.Sprintf("%s %s",msg.Chat.FirstName , msg.Chat.LastName)
	}
	return retVal
}

func(this TelegramBot)Send(msg tgbotapi.MessageConfig)(err error){
	var sentMsg tgbotapi.Message

	sentMsg, err = this.bot.Send(msg)
	if err == nil{
		this.log.Debug(fmt.Sprintf("Message %s sent successfuly to %s",sentMsg.Text, GetMsgContactTitle(sentMsg) ))
	}else{
		this.log.Error(err)
	}
	return err
}

func(this TelegramBot)SendFile(msg tgbotapi.DocumentConfig)(err error){
	var sentMsg tgbotapi.Message

	sentMsg, err = this.bot.Send(msg)
	if err == nil{
		this.log.Debug(fmt.Sprintf("Message %s sent successfuly to %s",sentMsg.Text, GetMsgContactTitle(sentMsg) ))
	}else{
		this.log.Error(err)
	}
	return err
}