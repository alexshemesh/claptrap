package telegram

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
)

func Test_NewTelegramBot(t *testing.T){
	bot := NewTelegramBot(*logs.NewLogger("test telegram bot"),"sometoken")
	if bot.token != "sometoken"{
		t.Errorf("Cannot initialize telegram bot")
	}
}