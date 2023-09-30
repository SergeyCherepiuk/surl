package validation

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func ValidateUrl(url string) error {
	if !urlRegexp.MatchString(url) {
		return fmt.Errorf("invalid url")
	}
	return nil
}

func ValidateExpiration(expiresIn int) error {
	if !slices.Contains(expiresInAllowedValues, expiresIn) {
		return fmt.Errorf("invalid expiration time")
	}
	return nil
}
