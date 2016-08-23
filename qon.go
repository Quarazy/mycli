package main

import (
	"os"

	"github.com/quarazy/qon/cmd"
)

func main() {
	commands := cmd.New()
	commands.Register("conf", cmd.GitConfigCmd)
	commands.Run(os.Args)
}
