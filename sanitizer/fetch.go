package sanitizer

import (
	"fmt"
	"regexp"
)

var urlRegex = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func SanitizeFetch(urls []string) error {
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
