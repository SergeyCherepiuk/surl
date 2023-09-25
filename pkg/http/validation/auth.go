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
	if !passwordMin8CharactersRegexp.MatchString(password) {
		return fmt.Errorf("password must be at least 8 characters long")
	} else if !passwordAllowedCharactersRegexp.MatchString(password) {
		return fmt.Errorf("only letters, numbers and special symbols are allowed")
	} else if !passwordMin1SpecialCharacterRegexp.MatchString(password) {
		return fmt.Errorf("password must contain at least 1 special character")
	} else if !passwordMin3UppercaseCharactersRegexp.MatchString(password) {
		return fmt.Errorf("password must contain at least 3 uppercase characters")
	}
	return nil
}
