package utils

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func UrlToFilepath(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	ret := parsedURL.Host + parsedURL.Path
	return strings.ReplaceAll(ret, string(filepath.Separator), "|")
}

func UrlToMetadataFilepath(urlStr string) string {
	filepath := UrlToFilepath(urlStr)
	return path.Join("."+filepath, "metadata.json")
}
