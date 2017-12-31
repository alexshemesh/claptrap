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
	"github.com/alexshemesh/claptrap/lib/contracts"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/markets/bitfinex"
	"fmt"
)


// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "get balance on bitfinex account",
	Long: `get balance on bitfinex account`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("bitfinex balance")
		vault := vault.NewVaultClientInitialized(log)
		vault.UserDomain = "alex_shemesh"
		err := runBitfinexBalance(log, vault)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	bitfinexCmd.AddCommand(balanceCmd)

}

func runBitfinexBalance( log logs.Logger, settings contracts.Settings)(err error){
	bitfinexClient := bitfinexClient.NewBitfinexClient(log, settings)
	var balance contracts.Balance
	balance,err = bitfinexClient.GetBalance()
	if err == nil {
		log.Log("Bitfinex Balances:")
		for _, pos := range(balance.Positions){
			log.Log(fmt.Sprintf("\t (%s)%s Amount %f Value %f", pos.Type, pos.Name, pos.Amount, pos.Value))
		}
	}
	return err
}
