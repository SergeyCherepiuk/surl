package validation

import (
	"fmt"
	"regexp"
)

// Source: https://stackoverflow.com/a/3809435
var urlRegexp = regexp.MustCompile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func ValidateOrigin(origin string) error {
	if !urlRegexp.MatchString(origin) {
		return fmt.Errorf("invalid url")
	}
	return nil
}