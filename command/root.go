package command

import (
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
)

var Command *cobra.Command

func init() {
	Command = &cobra.Command{
		Use: "fetch",
		RunE: func(cmd *cobra.Command, args []string) error {
			// testing purpose, to make sure we could write file
			filepath := time.Now().String()
			err := ioutil.WriteFile(filepath, []byte("Mantap betul"), 0755)
			return err
		},
	}
}
