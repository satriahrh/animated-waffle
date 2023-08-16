package main

import (
	"os"

	"github.com/satriahrh/autify-tht/command"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
