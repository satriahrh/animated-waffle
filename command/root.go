package command

import (
	"github.com/spf13/cobra"
)

var Command *cobra.Command

func init() {
	Command = &cobra.Command{
		Use: "fetch",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
}
