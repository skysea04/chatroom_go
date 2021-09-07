package test

import (
	"main/utils"
	"testing"
)

func TestValidEmail(t *testing.T) {
	a := utils.ValidEmail("arcade0425@gmail.com")
	c := utils.ValidEmail("123@.com")

	if a == true && c == false {
		t.Log("sucess")
	} else {
		t.Log("fail")
	}
}

func TestValidPwd(t *testing.T) {
	a := utils.ValidPwd("123asdfwqweqw")
	b := utils.ValidPwd("45sds")
	if a == true && b == false {
		t.Log("success")
	} else {
		t.Error("fail")
	}
}
