package validation

import "regexp"

var (
	usernameRegexp                        = regexp.MustCompile(`^[[:alnum:]_]{3,30}$`)
	passwordMin8CharactersRegexp          = regexp.MustCompile(`^.{8,}$`)
	passwordAllowedCharactersRegexp       = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*]+$`)
	passwordMin1SpecialCharacterRegexp    = regexp.MustCompile(`^(.*[!@#$%^&*]){1}.*$`)
	passwordMin3UppercaseCharactersRegexp = regexp.MustCompile(`^(.*[A-Z]){3}.*$`)

	// Source: https://stackoverflow.com/a/3809435
	urlRegexp              = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	expiresInAllowedValues = []int{1, 15, 30, 45, 60}
)
