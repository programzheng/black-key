package helper

import (
	"net/url"
)

func ValidateURL(str string) bool {
	parsedUrl, err := url.Parse(str)
	return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != ""
}
