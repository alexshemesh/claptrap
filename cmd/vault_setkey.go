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
	"github.com/spf13/viper"
	"github.com/alexshemesh/claptrap/lib/vault"
	"fmt"
)

var SecretPath string
var SecretValue string
// setkeyCmd represents the setkey command
var setkeyCmd = &cobra.Command{
	Use:   "setkey",
	Short: "Saves secret to vault",
	Long: `Saves secret to vault. Vault have to initialized and unsealed by vault init command`,
	Run: func(cmd *cobra.Command, args []string) {
		VaultAddr = viper.Get("vault.url").(string)
		VaultToken := viper.Get("vault.token").(string)
		log := *logs.NewLogger("vault setkey")
		err := runSetKey(log, VaultAddr,VaultToken, SecretPath, SecretValue)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},

}

func init() {
	vaultCmd.AddCommand(setkeyCmd)
	setkeyCmd.Flags().StringVarP(&SecretPath,"path" ,"p", "", "secret path")

	setkeyCmd.Flags().StringVarP(&SecretValue,"value" ,"v", "", "secret value")

}

func runSetKey(log logs.Logger, vaultAddr string, vaultToken string,secretPath string, secretValue string)(err error){
	vault := vault.NewVaultClient(vaultAddr,log)
	vault.Auth(vaultToken)
	err = vault.SetVaue(secretPath,secretValue)
	if err == nil {
		log.Log(fmt.Sprintf( "Success.Secret at path %s set. Value %s",secretPath, secretValue ))
	}
	return err
}