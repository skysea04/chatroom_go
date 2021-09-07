package utils

import (
	"net/mail"
	"regexp"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidPwd(pwd string) bool {
	pwdRule, _ := regexp.Compile("(^([a-zA-Z]+[0-9]+|[0-9]+[a-zA-Z]+)[a-zA-Z0-9]*$)")
	return pwdRule.MatchString(pwd) && len(pwd) >= 8
}
