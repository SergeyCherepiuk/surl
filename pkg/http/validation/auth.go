package validation

import (
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
)

func ValidateAuthentication(user domain.User) error {
	if !usernameRegexp.MatchString(user.Username) {
		return fmt.Errorf("invalid username")
	} else if !passwordRegexp.MatchString(user.Password) {
		return fmt.Errorf("invalid password")
	}
	return nil
}
