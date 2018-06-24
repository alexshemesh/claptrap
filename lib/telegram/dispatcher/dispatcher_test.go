package dispatcher

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/bouk/monkey"
	"strings"
)

func Test_NewDispatcher(t *testing.T){
	settingsPar := vault.NewVaultTestKit()
	disp := NewDispatcher(*logs.NewLogger("test dispatcher"),settingsPar)
	if disp == nil {
		t.Errorf("Unable to create telegram command dispatcher")
	}
}

func Test_UnknownCommand( t *testing.T){

	monkey.Patch( telegram.TelegramBot.Send, func(obj telegram.TelegramBot, msg tgbotapi.MessageConfig) (err error){
		return err
	})
	defer func(){
		monkey.Unpatch( telegram.TelegramBot.Send )
	}()
	settingsPar := vault.NewVaultTestKit()
	settingsPar.SetValue("telegram/token","sometoken")
	log := *logs.NewLogger("test dispatcher")
	disp := NewDispatcher(log,settingsPar)
	bot := telegram.NewTelegramBot(log, settingsPar)
	cmd := types.TGCommand{}
	cmd.Msg = &tgbotapi.Message{}
	cmd.Msg.MessageID = 24

	cmd.Msg.Chat = &tgbotapi.Chat{}
	cmd.Msg.Chat.ID = 99

	cmd.Msg.Text = "some bullshit"

	err := disp.Dispatch( *bot, cmd)
	if err != nil {
		t.Error(err)
	}
}

func Test_KunaOrdersCommand_NoAuth( t *testing.T){
	var sentMessage string

	monkey.Patch( telegram.TelegramBot.Send, func(obj telegram.TelegramBot, msg tgbotapi.MessageConfig) (err error){
		sentMessage = msg.Text
		return err
	})

	defer func(){
		monkey.Unpatch( telegram.TelegramBot.Send )
	}()

	settingsPar := vault.NewVaultTestKit()
	settingsPar.SetValue("telegram/token","sometoken")
	settingsPar.SetLogiResult(false)
	log := *logs.NewLogger("test dispatcher")
	disp := NewDispatcher(log,settingsPar)



	bot := telegram.NewTelegramBot(log, settingsPar)
	cmd := types.TGCommand{}
	cmd.Msg = &tgbotapi.Message{}
	cmd.Msg.MessageID = 24

	cmd.Msg.Chat = &tgbotapi.Chat{}
	cmd.Msg.Chat.ID = 99

	cmd.Msg.Text = "/kuna_orders_book"

	err := disp.Dispatch( *bot, cmd)
	if err != nil {
		t.Error(err)
	}

	 if strings.Index( sentMessage , "You are not authenticated to use system.\n Use /login <your password> command to log in into system" ) != 0 {
	 		t.Errorf("Should fail , user is not authenticated")
	 }
}

func Test_LoginCommand( t *testing.T){
	var sentMessage string

	monkey.Patch( telegram.TelegramBot.Send, func(obj telegram.TelegramBot, msg tgbotapi.MessageConfig) (err error){
		sentMessage = msg.Text
		return err
	})

	defer func(){
		monkey.Unpatch( telegram.TelegramBot.Send )
	}()

	settingsPar := vault.NewVaultTestKit()
	settingsPar.SetValue("telegram/token","sometoken")
	settingsPar.SetLogiResult(true)
	settingsPar.SetTokenToReturn("ABRACADABRA")
	log := *logs.NewLogger("test dispatcher")
	disp := NewDispatcher(log,settingsPar)



	bot := telegram.NewTelegramBot(log, settingsPar)
	cmd := types.TGCommand{}
	cmd.Msg = &tgbotapi.Message{}

	cmd.Msg.MessageID = 24

	cmd.Msg.Chat = &tgbotapi.Chat{}
	cmd.Msg.Chat.UserName = "vasya"
	cmd.Msg.Chat.ID = 99

	cmd.Msg.Text = "/login password"

	err := disp.Dispatch( *bot, cmd)
	if err != nil {
		t.Error(err)
	}


}