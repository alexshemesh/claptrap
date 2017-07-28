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
	"github.com/spf13/viper"
	"os"
	"github.com/alexshemesh/claptrap/lib/vault"
)
var UnsealKey string
// unsealCmd represents the unseal command
var unsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("vault unseal")
		VaultAddr = viper.Get("vault.url").(string)

		err := runUnseal(log, VaultAddr,UnsealKey)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	vaultCmd.AddCommand(unsealCmd)
	unsealCmd.Flags().StringVarP(&UnsealKey,"key" ,"k", "", "secret key to unseal the vault")

}

func runUnseal(log logs.Logger, vaultAddr string, key string)(err error){
	vaultClient := vault.NewVaultClient(vaultAddr,log)
	err = vaultClient.Unseal([]string{key})

	return err
}
