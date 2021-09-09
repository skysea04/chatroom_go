package views

import (
	"main/db_client"
	"strconv"

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
	var roomID int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(302, "/")
	}
	err = db_client.DB.QueryRow("SELECT id FROM rooms WHERE id = ?;", id).Scan(&roomID)
	if err != nil {
		c.Redirect(302, "/")
	}
	return c.Render(200, "chatroom.html", nil)
}
