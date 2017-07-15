package main

import "github.com/alexshemesh/claptrap/cmd"
var Version string

func main() {
	print("Version:" + Version + "\n")
	cmd.ProgramVersion = Version
	cmd.Execute()
}
