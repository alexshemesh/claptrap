// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/vault"
	"os"
	"github.com/alexshemesh/claptrap/lib/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

// checkmsgCmd represents the checkmsg command
var checkmsgCmd = &cobra.Command{
	Use:   "checkmsg",
	Short: "Check for incoming messages",
	Long: `Will check for incoming messages and answer with simple
	text.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("telegram checkmsg")

		vaultClient := vault.NewVaultClientInitialized(log)
		token,err := vaultClient.GetValue("telegram/token")
		if err == nil {
			err := runCheckMsg(log, token)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
		}
		os.Exit(0 )
	},
}

func init() {
	telegramCmd.AddCommand(checkmsgCmd)
}

func procMessage(log logs.Logger,bot telegram.TelegramBot, msg tgbotapi.Message)(err error){
	log.Log(fmt.Sprintf( "User %s sent us message %s", msg.From.UserName, msg.Text))

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "Got your message")
	newMsg.ReplyToMessageID = msg.MessageID

	err = bot.Send(newMsg)
	return err
}

func runCheckMsg(log logs.Logger, token string)(err error){
	telegramBot := telegram.NewTelegramBot(log,token)
	err = telegramBot.Connect()
	if err == nil{
		var messages []tgbotapi.Message

		messages,err = telegramBot.GetMessages()
		log.Log(fmt.Sprintf("Got %d messages", len(messages)))
		for i, msg := range(messages){
			log.Log(fmt.Sprintf("Processing message %d", i + 1))
			err = procMessage(log,*telegramBot,msg)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return err
}