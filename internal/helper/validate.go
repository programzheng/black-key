package helper

import (
	"net/url"
)

func ValidateURL(str string) error {
	_, err := url.Parse(str)
	return err
}
