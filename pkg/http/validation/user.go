package validation

import (
	"fmt"
	"regexp"

	"github.com/SergeyCherepiuk/surl/domain"
)

var usernameRegexp = regexp.MustCompile("^[[:alnum:]_]{3, 30}$")
var passwordRegexp = regexp.MustCompile("^[[:alnum:]@$!%*#?&]{8,}$")

func ValidateUserSignUp(user domain.User) error {
	if (!usernameRegexp.MatchString(user.Username)) {
		return fmt.Errorf("invalid username")
	} else if (!passwordRegexp.MatchString(user.Password)) {
		return fmt.Errorf("invalid password")
	}
	return nil 
}
