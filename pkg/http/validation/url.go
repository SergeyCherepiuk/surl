package validation

import (
	"fmt"
)

func ValidateUrl(url string) error {
	if !urlRegexp.MatchString(url) {
		return fmt.Errorf("invalid url")
	}
	return nil
}
