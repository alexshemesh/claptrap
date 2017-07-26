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

func (this TelegramBot) Ping()(err error){
	if this.bot.Self.UserName == "" {
		err = fmt.Errorf("Bot is not connected")
	}
	this.log.Debug("Telegram Ping is OK. Connected with username " + this.bot.Self.UserName)
	return err
}