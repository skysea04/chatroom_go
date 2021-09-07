package views

import (
	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	return c.Render(200, "signup.html", nil)
}

func Index(c echo.Context) error {
	return c.Render(200, "index.html", nil)
}

func MyRooms(c echo.Context) error {
	return c.Render(200, "my_rooms.html", nil)
}

func CreateRoom(c echo.Context) error {
	return c.Render(200, "create_room.html", nil)
}

func ChatRoom(c echo.Context) error {
	return c.Render(200, "chatroom.html", nil)
}
