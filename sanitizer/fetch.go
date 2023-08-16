package sanitizer

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var urlRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

var SanitizeFetch cobra.PositionalArgs = func(_ *cobra.Command, urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("no url given")
	}

	for _, url := range urls {
		if !urlRegex.MatchString(url) {
			return fmt.Errorf("invalid url")
		}
	}
	return nil
}
