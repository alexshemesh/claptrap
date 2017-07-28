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
	"github.com/spf13/viper"
)


// issealedCmd represents the issealed command
var issealedCmd = &cobra.Command{
	Use:   "issealed",
	Short: "check if vault is sealed",
	Long: `Whenever vault server is restarted it becomes sealed
	and have to be unsealed before usage. Unseal procedure takes a key/s to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("vault issealed")
		VaultAddr = viper.Get("vault.url").(string)

		err := runIsSealed(log, VaultAddr)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	vaultCmd.AddCommand(issealedCmd)



}

func runIsSealed(log logs.Logger,vaultAddr string)(err error ){

	vault := vault.NewVaultClient(vaultAddr,log)
	var sealed bool
	sealed,err = vault.IsSealed()
	if sealed {
		log.Log("Vault is sealed")
	}else{
		log.Log("Vault is NOT sealed. Ready to use.")
	}
	return err
}