package validation

import (
	"fmt"
)

func ValidateOrigin(origin string) error {
	if !urlRegexp.MatchString(origin) {
		return fmt.Errorf("invalid url")
	}
	return nil
}
