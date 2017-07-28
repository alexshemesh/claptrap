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
	"github.com/spf13/viper"
	"os"
)


// getkeyCmd represents the getkey command
var getkeyCmd = &cobra.Command{
	Use:   "getkey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		VaultAddr = viper.Get("vault.url").(string)
		VaultToken := viper.Get("vault.token").(string)
		log := *logs.NewLogger("vault getkey")
		err := runGetKey(log, VaultAddr,VaultToken, SecretPath)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	vaultCmd.AddCommand(getkeyCmd)
	getkeyCmd.Flags().StringVarP(&SecretPath,"path" ,"p", "", "secret path")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func runGetKey(log logs.Logger, vaultAddr string, vaultToken string,secretPath string)(err error){
	vault := vault.NewVaultClient(vaultAddr,log)
	vault.Auth(vaultToken)
	var secretValue string
	secretValue,err = vault.GetValue(secretPath)
	if err == nil {
		log.Log(fmt.Sprintf( "Success.Secret at path %s. Value %s",secretPath, secretValue ))
	}
	return err
}