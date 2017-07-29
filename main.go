package main

import (
	"github.com/alexshemesh/claptrap/cmd"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
)
var Version string

func main() {
	homeFolder := cmd.UserHomeDir()
	content,_ := ioutil.ReadFile(path.Join(homeFolder,"claptrap"))
	if content != nil{
		print("Some data\n" )
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName(".claptrap")
	viper.AddConfigPath(homeFolder)
	err := viper.ReadInConfig()
	if err != nil {
		print("No configuration file found\n" + err.Error() + "\n")
	}
	print("Version:" + Version + "\n")
	cmd.ProgramVersion = Version
	cmd.Execute()
}
