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
	"net/url"
	"path"
	"github.com/alexshemesh/claptrap/lib/http"
	"fmt"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init vault",
	Long: `Initialize fresh vault instance,
	if vault is already intialized - nothing should happen.
	Confiuration information will be saved into Viper configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		log := *logs.NewLogger("vault")
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

func runInit( log logs.Logger , vaultAddr string) ( err error ){
	http := httpClient.NewHttpExecutor()
	body :=`{"secret_shares":1,"secret_threshold":1}`

	var u *url.URL
	u, err = url.Parse(vaultAddr)
	u.Path = path.Join(u.Path, "v1", "sys", "init")
	var response []byte
	response, err = http.Get().Execute( u.String(), nil, []byte(body) )
	if err == nil {
		if response == nil {
			err = fmt.Errorf("No response")
			log.Error(err)
		}else{
			log.Log("Vault initialized")
		}
	}

	return err
}