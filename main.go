package main

import (
	"os"

	"github.com/satriahrh/animated-waffle/command"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
