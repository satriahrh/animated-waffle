package main

import (
	"fmt"
	"os"

	"github.com/satriahrh/autify-tht/command"
)

func main() {
	if err := command.Command.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
