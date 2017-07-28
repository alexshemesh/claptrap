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


	"github.com/spf13/cobra"
	"github.com/alexshemesh/claptrap/lib/logs"
	"os"

	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/telegram"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "pings telegram",
	Long: `Checks that telegram bot is configured
	properly and it can connect to telegram`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("telegram ping")

		vaultClient := vault.NewVaultClientInitialized(log)
		token,err := vaultClient.GetValue("telegram/token")
		if err == nil {
			err := runPing(log, token)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
		}
		os.Exit(0 )
	},
}

func init() {
	telegramCmd.AddCommand(pingCmd)

}

func runPing(log logs.Logger, token string)(err error){
	telegramBot := telegram.NewTelegramBot(log,token)
	err = telegramBot.Connect()
	if err == nil{
		err = telegramBot.Ping()
	}
	if err == nil {
		log.Log("Telegram ping returns OK")
	}
	return err
}