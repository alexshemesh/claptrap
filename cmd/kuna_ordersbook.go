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
	"github.com/alexshemesh/claptrap/lib/markets/kuna"
)

// ordersbookCmd represents the ordersbook command
var ordersbookCmd = &cobra.Command{
	Use:   "ordersbook",
	Short: "Retruns book of orders frm Kuna market",
	Long: `Retruns book of orders frm Kuna market`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("kuna ordersbook")
		err := runKunaOrdersBook(log)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	kunaCmd.AddCommand(ordersbookCmd)



}

func runKunaOrdersBook(log logs.Logger)(err error){

	kunaClient := kuna.NewKunaClient(log)

	book,err := kunaClient.GetOrdersBook()
	if err == nil {
		log.Log("Asks")
		for i, order := range (book.Asks){
			log.Log(log.CS.Sprintf("\t%d %s %s %s %s", i, order.Side, order.OrdType, log.CS.Green(order.Price),  order.Volume ) )
		}
		log.Log("Bids")
		for i, order := range (book.Bids) {
			log.Log(log.CS.Sprintf("\t%d %s %s %s %s", i, order.Side, order.OrdType, log.CS.Red(order.Price), order.Volume))
		}
	}
	return err
}