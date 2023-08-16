package command

import (
	"github.com/satriahrh/autify-tht/sanitizer"
	"github.com/spf13/cobra"
)

var Command *cobra.Command

func init() {
	Command = &cobra.Command{
		Use:  "fetch {set of urls}",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), sanitizer.SanitizeFetch),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
