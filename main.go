package main

import (
	"github.com/bizshuk/cc-plugin/cmd"
	"github.com/bizshuk/cc-plugin/config"
)

func main() {
	config.Init()
	cmd.Execute()
}
