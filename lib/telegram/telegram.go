package telegram

import (
	"github.com/alexshemesh/claptrap/lib/logs"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"

)

type TelegramBot struct{
	log logs.Logger
	token string
	bot *tgbotapi.BotAPI
}


func NewTelegramBot(logPar logs.Logger, tokenPar string)(retVal *TelegramBot) {
	retVal = &TelegramBot{log: logPar, token: tokenPar }

	return retVal
}

func (this *TelegramBot)Connect()(err error){
	this.bot, err = tgbotapi.NewBotAPI(this.token)
	return err
}

func (this TelegramBot)GetMessages( )(retVal []tgbotapi.Message,err error) {
	u := tgbotapi.NewUpdate(0)
	var updates []tgbotapi.Update
	updates, err = this.bot.GetUpdates(u)
	if err == nil {
		for _, update := range updates {
			if update.Message != nil {
				retVal = append(retVal, *update.Message)
			}
		}
	}
	return retVal,err
}

func (this TelegramBot) Ping()(err error){
	if this.bot.Self.UserName == "" {
		err = fmt.Errorf("Bot is not connected")
	}
	this.log.Debug("Telegram Ping is OK. Connected with username " + this.bot.Self.UserName)
	return err
}

func(this TelegramBot)Send(msg tgbotapi.MessageConfig)(err error){
	var sentMsg tgbotapi.Message
	sentMsg, err = this.bot.Send(msg)
	if err == nil{
		this.log.Debug(fmt.Sprintf("Message %s sent successfuly to %s %s",sentMsg.Text, sentMsg.Contact.FirstName , sentMsg.Contact.LastName))
	}else{
		this.log.Error(err)
	}
	return err
}