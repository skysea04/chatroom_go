package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/db_client"
	"main/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	RePwd string `json:"rePwd"`
	Name  string `json:"name"`
}

type ErrMsg struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

var JwtSecret = []byte("secret")

// 註冊
func PostUser(c echo.Context) error {
	var reqBody User
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reqBody)
	if err != nil {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "輸入格式錯誤",
		})
	}

	// 檢查email是否符合格式
	mailIsValid := utils.ValidEmail(reqBody.Email)
	if !mailIsValid {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "信箱不符合格式，請重新輸入",
		})
	}

	// 檢查密碼是否符合格式
	pwdIsValid := utils.ValidPwd(reqBody.Pwd)
	if !pwdIsValid {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "請輸入8位以上的英數組合",
		})
	}

	// 檢查密碼與重新輸入的密碼是否相符
	if reqBody.Pwd != reqBody.RePwd {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "密碼與重新輸入密碼不相符",
		})
	}

	// 檢查暱稱是否符合格式
	if len(reqBody.Name) < 4 {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "暱稱最少為4字元",
		})
	}

	var user User
	// 檢查帳號是否重複
	err = db_client.DBClient.QueryRow("SELECT email FROM users WHERE email = ?;", reqBody.Email).Scan(&user.Email)
	if err == nil {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "該信箱已註冊",
		})
	}

	// 檢查暱稱是否重複
	err = db_client.DBClient.QueryRow("SELECT name FROM users WHERE name = ?;", reqBody.Name).Scan(&user.Name)
	if err == nil {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "該暱稱已被使用",
		})
	}

	_, err = db_client.DBClient.Exec("INSERT INTO users (email, password, name) VALUES (?, ?, ?);", reqBody.Email, reqBody.Pwd, reqBody.Name)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   fmt.Sprint(err),
		})
	}
	return c.JSON(200, echo.Map{
		"ok":  true,
		"msg": "註冊成功！",
	})
}

// 登入
func PatchUser(c echo.Context) error {
	var reqBody User
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reqBody)
	if err != nil {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "輸入格式錯誤",
		})
	}

	// 確認使用者存在
	var user User
	err = db_client.DBClient.QueryRow("SELECT id, name FROM users WHERE email = ? AND password = ?;", reqBody.Email, reqBody.Pwd).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(401, ErrMsg{
				Error: true,
				Msg:   "帳號或密碼輸入錯誤",
			})
		}
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   err.Error(),
		})
	}

	// 建一個 token 給 user
	claims := &JwtCustomClaims{
		user.ID,
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(JwtSecret)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	// 回傳token
	return c.JSON(200, echo.Map{
		"ok": true,
	})
}

func GetUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	name := claims.Name
	id := claims.ID
	return c.JSON(200, echo.Map{
		"ok":   true,
		"id":   id,
		"name": name,
	})
}

func JwtGateKeeper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.Redirect(302, "/signup")
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &JwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
			c.Set("userID", claims.ID)
			c.Set("userName", claims.Name)
			// log.Printf("%v, %v", claims.ID, claims.Name)
			return next(c)
		} else {
			return c.Redirect(302, "/signup")
		}
	}
}
