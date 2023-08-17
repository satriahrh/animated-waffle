package sanitizer

import (
	"fmt"

	"github.com/satriahrh/autify-tht/utils"
	"github.com/spf13/cobra"
)

var SanitizeFetch cobra.PositionalArgs = func(_ *cobra.Command, urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("no url given")
	}

	for _, url := range urls {
		if !utils.UrlRegex.MatchString(url) {
			return fmt.Errorf("invalid url")
		}
	}
	return nil
}
