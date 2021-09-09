package utils

import (
	"crypto/sha512"
	"fmt"
)

func HashPassword(email string, pwd string) string {
	hashedPwd := fmt.Sprintf("%x", sha512.Sum512([]byte(email+pwd)))
	return hashedPwd
}

func VerifyHash(email string, pwd string, hashedPwd string) bool {
	pwdNeedVerify := fmt.Sprintf("%x", sha512.Sum512([]byte(email+pwd)))
	return pwdNeedVerify == hashedPwd
}
