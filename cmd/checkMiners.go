// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"os"

	"github.com/spf13/cobra"
	"github.com/alexshemesh/claptrap/lib/logs"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/types"
	"github.com/alexshemesh/claptrap/lib/claymore"
	"github.com/alexshemesh/claptrap/lib/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"time"
)

// checkMinersCmd represents the checkMiners command
var checkMinersCmd = &cobra.Command{
	Use:   "checkMiners",
	Short: "Check miners status",
	Long: `Check miners status. Send notification to telegram if something changes`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("telegram checkmsg")

		vaultClient := vault.NewVaultClientInitialized(log)
		err := runLoop(log, vaultClient,vaultClient)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	RootCmd.AddCommand(checkMinersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkMinersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkMinersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runMinersCheck(log logs.Logger,settings types.Settings,auth types.Auth)(err error){
	telegramBot := telegram.NewTelegramBot(log,settings)
	err = telegramBot.Connect()

	claymoreClient := claymore.NewClaymoreManagerClient(log, settings)
	res,reasons, err := claymoreClient.CheckAndCompare()
	//-1001344868791
	//-1001187769131 - Poti ops
	var newMsg tgbotapi.MessageConfig
	if err != nil {
		newMsg = tgbotapi.NewMessage(-1001187769131, err.Error())
	}else if res != true{
		newMsg = tgbotapi.NewMessage(-1001187769131, strings.Join(reasons, "\n"))
	}

	if newMsg.ChatID != 0 {
		err = telegramBot.Send(newMsg)
	}

	return err
}


func runLoop(log logs.Logger,settings types.Settings,auth types.Auth)(err error){
	for true {
		runMinersCheck(log,settings,auth)
		time.Sleep(time.Second * 30)
	}

	return err
}

