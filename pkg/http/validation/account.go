package validation

import "fmt"

func ValidateUsernameChange(username string) error {
	if !usernameRegexp.MatchString(username) {
		return fmt.Errorf("invalid username")
	}
	return nil
}
