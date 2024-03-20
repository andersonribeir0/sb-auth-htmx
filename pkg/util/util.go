package util

import "regexp"

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	return emailRegex.MatchString(email)
}

func ValidateStrongPassword(password string) bool {
	var (
		minLengthRegex       = regexp.MustCompile(`.{8,}`)
		uppercaseLetterRegex = regexp.MustCompile(`[A-Z]`)
		lowercaseLetterRegex = regexp.MustCompile(`[a-z]`)
		numberRegex          = regexp.MustCompile(`[0-9]`)
		specialCharRegex     = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	)

	return minLengthRegex.MatchString(password) &&
		uppercaseLetterRegex.MatchString(password) &&
		lowercaseLetterRegex.MatchString(password) &&
		numberRegex.MatchString(password) &&
		specialCharRegex.MatchString(password)
}
