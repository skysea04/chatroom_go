package utils

import (
	"net/mail"
	"regexp"
)

var pwdRule, _ = regexp.Compile("[a-zA-Z0-9]")

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidPwd(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}
	return pwdRule.MatchString(pwd)
}
