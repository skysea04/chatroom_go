package utils

import (
	"log"

	"github.com/matthewhartstonge/argon2"
)

var argon = argon2.DefaultConfig()

func HashPassword(pwd string) (hashedPwd string, err error) {
	encodedPwd, err := argon.HashEncoded([]byte(pwd))
	if err != nil {
		log.Println(err)
		return
	}
	hashedPwd = string(encodedPwd)
	return
}

func VerifyHash(pwd string, hashedPwd string) bool {
	ok, err := argon2.VerifyEncoded([]byte(pwd), []byte(hashedPwd))
	if err != nil {
		log.Println(err)
	}
	return ok
}
