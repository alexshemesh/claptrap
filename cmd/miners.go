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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/alexshemesh/claptrap/lib/http"
	"net/url"
	"encoding/json"
	"path"
	"net"
	"bufio"
	"github.com/alexshemesh/claptrap/lib/claymore"
)

// minersCmd represents the miners command
var minersCmd = &cobra.Command{
	Use:   "miners",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("miners called")
		//GetStatisticsTCP()
			GetMinerStatistic()
	},
}

func init() {
	RootCmd.AddCommand(minersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// minersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// minersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetMinerStatistic() (retVal string, err error) {

	httpClient := httpClient.NewHttpExecutor().WithBasicAuth("admin", "statuscheck")
	var u *url.URL
	u, err = url.Parse(path.Join(""))
	u.Scheme = "http"
	u.Host = "10.7.7.2:8193"

	q := u.Query()
	u.RawQuery = q.Encode()
	var response []byte
	body := `{"id":0,"jsonrpc":"2.0","method":"miner_getstat"}`
	response, err = httpClient.Post().Execute(u.String(), nil, []byte(body))
	retVal = string(response)
	if err == nil {
		err = json.Unmarshal(response, &retVal)
	}
	print(retVal)
	claymore.SplitTable(string(response))
	return retVal, err
}

func GetStatisticsTCP() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "10.7.7.2:8193")

	buffer := make([]byte, 2048)
	// send to socket
	fmt.Fprintln(conn, `{"id":0,"jsonrpc":"2.0","method":"miner_getstat1","psw":"statuscheck"}`)
	// listen for reply
	reader := bufio.NewReader(conn)
	var err error
	var b byte
	for err == nil {
		b, err = reader.ReadByte()
		print(b)
	}

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {

		fmt.Print("Message from server: " + string(buffer))
		claymore.SplitTable(string(buffer))
	}
}
