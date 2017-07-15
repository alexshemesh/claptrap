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

)

var VaultAddr string
// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Set of commands to manage vault",
	Long: `Hashicorp Vault is used for storing sensitive information
	account details, configurations etc.
	This set of commands allows to initialize new vault, change and monitor values in vault`,
	Run: func(cmd *cobra.Command, args []string ) {

	},
}

func init() {
	RootCmd.AddCommand(vaultCmd)
	vaultCmd.PersistentFlags().StringVarP(&VaultAddr,"vaultAddr" ,"", "http://127.0.0.1", "address of vault server")

}

