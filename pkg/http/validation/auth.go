package validation

import (
	"fmt"
)

func ValidateUsername(username string) error {
	if !usernameRegexp.MatchString(username) {
		return fmt.Errorf("invalid username")
	}
	return nil
}

func ValidatePassword(password string) error {
	if !passwordRegexp.MatchString(password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
