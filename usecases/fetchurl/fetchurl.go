package fetchurl

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/satriahrh/autify-tht/adapters"
)

func Construct(fetchContent adapters.FetchContent, storeContent adapters.StoreContent) func(ctx context.Context, url string) error {
	return func(ctx context.Context, url string) error {
		reader, err := fetchContent(ctx, url)
		if err != nil {
			return err
		}

		var filepath string = urlToFilepath(url)
		filepath += ".html"

		err = storeContent(ctx, filepath, reader)
		if err != nil {
			return err
		}

		return nil
	}
}

func urlToFilepath(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	ret := parsedURL.Host + parsedURL.Path
	return strings.ReplaceAll(ret, string(filepath.Separator), "|")
}
