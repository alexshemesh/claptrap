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
		"fmt"
	"github.com/alexshemesh/claptrap/lib/vault"
	"runtime"
)

type VaultSettings struct {
	VaultURL string `json:"vaulturl"`
	Secret string `json:"secret"`
	Token string `json:"token"`
}
// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init vault",
	Long: `Initialize fresh vault instance,
	if vault is already intialized - nothing should happen.
	Confiuration information will be saved into Viper configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("vault init")
		err := runInit(log, VaultAddr)
		if err != nil {
			log.Error(err)
			os.Exit(-1)
		}
		os.Exit(0 )
	},
}

func init() {
	vaultCmd.AddCommand(initCmd)
}

func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

func runInit( log logs.Logger , vaultAddr string) ( err error ){

	vault := vault.NewVaultClient(vaultAddr,*logs.NewLogger("root"))
	var initialized bool
	initialized,err = vault.IsInitialized()
	if initialized == false {
		keys,token,err := vault.Initialize()
		if err == nil {
			vault.Unseal(keys)

			log.Log(fmt.Sprintf("Vault at % initialized and unsealed", vaultAddr))
			log.Log("Save these keys in safe location. You will not be able to retrieve them from server againg!!!")
			for i, key := range (keys) {
				log.Log(fmt.Sprintf("\tKey %d: %s", i, key))
			}
			log.Log(fmt.Sprintf("Add following lines to your %s/.clpatrap.yaml configuration file\nvault:\n  url: %s\n  token: %s\n", UserHomeDir(),vaultAddr,token))

		}
	} else {
		err = fmt.Errorf("Already initialized")
	}
	return err
}