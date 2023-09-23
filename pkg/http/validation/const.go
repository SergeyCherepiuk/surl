package validation

import "regexp"

var (
	usernameRegexp = regexp.MustCompile(`^[[:alnum:]_]{3,30}$`)
	passwordRegexp = regexp.MustCompile(`^[[:alnum:]@$!%*#?&]{8,}$`) // TODO: Improve validation to  require uppercase letters and digits

	// Source: https://stackoverflow.com/a/3809435
	urlRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
)
