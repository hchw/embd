package main

import (
	"os"

	_ "github.com/kidoman/embd/host/all"
	"github.com/urfave/cli"
)

var version = "0.1.0"

var commands []cli.Command

func registerCommand(cmd cli.Command) {
	commands = append(commands, cmd)
}

func main() {
	app := cli.NewApp()
	app.Name = "embd"
	app.Usage = "embedded utility belt"
	app.Version = version
	app.Commands = commands

	app.Run(os.Args)
}
