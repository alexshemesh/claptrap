package telegram

import (
	"github.com/alexshemesh/claptrap/lib/contracts"
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/bouk/monkey"
	"gopkg.in/telegram-bot-api.v4"
	"encoding/json"
	"fmt"
)

func Test_NewTelegramBot(t *testing.T){
	settings := vault.NewVaultTestKit()
	settings.SetValue("telegram/token","sometoken")
	bot := NewTelegramBot(*logs.NewLogger("test telegram bot"),settings)
	if bot.token != "sometoken"{
		t.Errorf("Cannot initialize telegram bot")
	}
}

func Test_GetMessages_NoPreviousData(t *testing.T){
	monkey.Patch((*tgbotapi.BotAPI).GetUpdates, func(this *tgbotapi.BotAPI,config tgbotapi.UpdateConfig)(retVal []tgbotapi.Update,err error){
		message1 := tgbotapi.Update{}
		message1.UpdateID = 15
		message1.Message = &tgbotapi.Message{}
		retVal = append(retVal, message1)

		message2 := tgbotapi.Update{}
		message2.UpdateID = 16
		message2.Message = &tgbotapi.Message{}
		retVal = append(retVal, message2)

		return 	retVal,err
	} )

	monkey.Patch( tgbotapi.NewBotAPI , func(token string)(retVal *tgbotapi.BotAPI,err error){
		retVal = &tgbotapi.BotAPI{}
		return retVal,err
	})

	defer func (){
		monkey.Unpatch((*tgbotapi.BotAPI).GetUpdates)
		monkey.Unpatch( tgbotapi.NewBotAPI )
	}()

	settings := vault.NewVaultTestKit()
	settings.SetValue("telegram/token","sometoken")

	bot := NewTelegramBot(*logs.NewLogger("test telegram bot"),settings)
	bot.Connect()
	messages,err := bot.GetMessages()
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 2 {
		t.Error("Cannot receive all the messages")
	}
	var newSettings string
	newSettings,err = settings.GetValue("telegram/settings")

	if newSettings == "" {
		t.Error("No telegram data")
	}

	var newSettingsStruct contracts.TelegramSettings
	err = json.Unmarshal([]byte(newSettings), &newSettingsStruct )
	if err != nil {
		t.Error(err)
	}
	if newSettingsStruct.LastUpdate != 16 {
		t.Error("Settings where not set properly")
	} 
}

func Test_GetMessages_WithPreviousData(t *testing.T){
	currentOffset := 23
	monkey.Patch((*tgbotapi.BotAPI).GetUpdates, func(this *tgbotapi.BotAPI,config tgbotapi.UpdateConfig)(retVal []tgbotapi.Update,err error) {
			if config.Offset == currentOffset+1 {
				message1 := tgbotapi.Update{}
				message1.UpdateID = 23
				message1.Message = &tgbotapi.Message{}
				retVal = append(retVal, message1)
			}
			return retVal, err
		})

	monkey.Patch( tgbotapi.NewBotAPI , func(token string)(retVal *tgbotapi.BotAPI,err error){
		retVal = &tgbotapi.BotAPI{}
		return retVal,err
	})

	defer func (){
		monkey.Unpatch((*tgbotapi.BotAPI).GetUpdates)
		monkey.Unpatch( tgbotapi.NewBotAPI )
	}()

	settings := vault.NewVaultTestKit()
	settings.SetValue("telegram/token","sometoken")
	settings.SetValue("telegram/settings",fmt.Sprintf(`{"LastUpdate":%d}`, currentOffset) )

	bot := NewTelegramBot(*logs.NewLogger("test telegram bot"),settings)
	bot.Connect()
	messages,err := bot.GetMessages()
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 1 {
		t.Error("Cannot receive all the messages")
	}
	var newSettings string
	newSettings,err = settings.GetValue("telegram/settings")

	if newSettings == "" {
		t.Error("No telegram data")
	}

	var newSettingsStruct contracts.TelegramSettings
	err = json.Unmarshal([]byte(newSettings), &newSettingsStruct )
	if err != nil {
		t.Error(err)
	}
	if newSettingsStruct.LastUpdate != 23 {
		t.Error("Settings where not set properly")
	}
}

func Test_GetMessages_404(t *testing.T){
	currentOffset := 0
	monkey.Patch((*tgbotapi.BotAPI).GetUpdates, func(this *tgbotapi.BotAPI,config tgbotapi.UpdateConfig)(retVal []tgbotapi.Update,err error) {
		if config.Offset == currentOffset+1 {
			message1 := tgbotapi.Update{}
			message1.UpdateID = 23
			message1.Message = &tgbotapi.Message{}
			retVal = append(retVal, message1)
		}
		return retVal, err
	})

	monkey.Patch( tgbotapi.NewBotAPI , func(token string)(retVal *tgbotapi.BotAPI,err error){
		retVal = &tgbotapi.BotAPI{}
		return retVal,err
	})

	defer func (){
		monkey.Unpatch((*tgbotapi.BotAPI).GetUpdates)
		monkey.Unpatch( tgbotapi.NewBotAPI )
	}()

	settings := vault.NewVaultTestKit()
	settings.SetValue("telegram/token","sometoken")
	settings.SetErrorEmulation("telegram/settings", fmt.Errorf("404 Not Found"))

	bot := NewTelegramBot(*logs.NewLogger("test telegram bot"),settings)
	bot.Connect()
	messages,err := bot.GetMessages()
	if err != nil {
		t.Error(err)
	}
	if len(messages) != 1 {
		t.Error("Cannot receive all the messages")
	}
	var newSettings string
	newSettings,err = settings.GetValue("telegram/settings")

	if newSettings == "" {
		t.Error("No telegram data")
	}

	var newSettingsStruct contracts.TelegramSettings
	err = json.Unmarshal([]byte(newSettings), &newSettingsStruct )
	if err != nil {
		t.Error(err)
	}
	if newSettingsStruct.LastUpdate != 23 {
		t.Error("Settings where not set properly")
	}
}

