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
	viper.SetConfigName(".claptrap") // name of config file (without extension)
	viper.AddConfigPath(homeFolder)  // call multiple times to add many search paths
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		print("No configuration file found\n" + err.Error() + "\n")
	}
	print("Version:" + Version + "\n")
	cmd.ProgramVersion = Version
	cmd.Execute()
}
