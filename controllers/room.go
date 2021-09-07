package controllers

import (
	"encoding/json"
	"fmt"
	"main/db_client"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Room struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	CreatedAt string `json:"createdAt"`
	Url       string `json:"url"`
}

type PageStatus struct {
	Page    int `json:"page"`
	MaxPage int `json:"maxPage"`
}

func PostRoom(c echo.Context) error {

	var reqBody Room
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reqBody)
	if err != nil {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "輸入格式錯誤",
		})
	}
	if len(reqBody.Name) < 4 {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "房間名稱最少為4字元",
		})
	}

	userID := c.Get("userID")
	_, err = db_client.DB.Exec("INSERT INTO rooms(name, owner) VALUES (?, ?);", reqBody.Name, userID)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}

	return c.JSON(200, echo.Map{
		"ok":  true,
		"msg": "建立成功!",
	})
}

// 獲取所有聊天室第n頁的資料
func GetRooms(c echo.Context) error {
	var roomCount int
	var pageStatus PageStatus
	var rooms []Room

	if c.QueryParam("page") == "" {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "輸入格式錯誤",
		})
	}
	pageStatus.Page, _ = strconv.Atoi(c.QueryParam("page"))

	// 獲取聊天室總數
	rows, err := db_client.DB.Query("SELECT COUNT(id) FROM rooms;")
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&roomCount); err != nil {
			return c.JSON(500, ErrMsg{
				Error: true,
				Msg:   "伺服器內部錯誤",
			})
		}
	}

	pageStatus.MaxPage = int(math.Ceil(float64(roomCount) / 10))

	if pageStatus.Page > pageStatus.MaxPage && pageStatus.MaxPage > 0 {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "頁數不符",
		})
	}

	// 獲取聊天室資料
	rows, err = db_client.DB.Query("SELECT rooms.id, rooms.name, users.name FROM rooms INNER JOIN users ON rooms.owner = users.id ORDER BY rooms.id DESC LIMIT ?, 10;", (pageStatus.Page-1)*10)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}
	defer rows.Close()
	for rows.Next() {
		var singleRoom Room
		if err := rows.Scan(&singleRoom.ID, &singleRoom.Name, &singleRoom.Owner); err != nil {
			return c.JSON(500, ErrMsg{
				Error: true,
				Msg:   "伺服器內部錯誤",
			})
		}
		singleRoom.Url = fmt.Sprintf("/chatroom/%v", singleRoom.ID)
		rooms = append(rooms, singleRoom)
	}

	return c.JSON(200, echo.Map{
		"pageStatus": pageStatus,
		"rooms":      rooms,
	})
}

func GetMyRooms(c echo.Context) error {
	var roomCount int
	var pageStatus PageStatus
	var rooms []Room
	userID := c.Get("userID")

	if c.QueryParam("page") == "" {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "輸入格式錯誤",
		})
	}
	pageStatus.Page, _ = strconv.Atoi(c.QueryParam("page"))

	// 獲取個人擁有的聊天室總數
	err := db_client.DB.QueryRow("SELECT COUNT(owner) FROM rooms WHERE owner = ?;", userID).Scan(&roomCount)

	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}

	pageStatus.MaxPage = int(math.Ceil(float64(roomCount) / 10))

	if pageStatus.Page > pageStatus.MaxPage && pageStatus.MaxPage > 0 {
		return c.JSON(400, ErrMsg{
			Error: true,
			Msg:   "頁數不符",
		})
	}

	// 獲取個人聊天室資料
	rows, err := db_client.DB.Query("SELECT rooms.id, rooms.name, users.name FROM rooms INNER JOIN users ON rooms.owner = users.id WHERE users.id = ? ORDER BY rooms.id DESC LIMIT ?, 10;", userID, (pageStatus.Page-1)*10)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}
	defer rows.Close()
	for rows.Next() {
		var singleRoom Room
		if err := rows.Scan(&singleRoom.ID, &singleRoom.Name, &singleRoom.Owner); err != nil {
			return c.JSON(500, ErrMsg{
				Error: true,
				Msg:   "伺服器內部錯誤",
			})
		}
		singleRoom.Url = fmt.Sprintf("/chatroom/%v", singleRoom.ID)
		rooms = append(rooms, singleRoom)
	}

	return c.JSON(200, echo.Map{
		"pageStatus": pageStatus,
		"rooms":      rooms,
	})
}

func GetRoomInfo(c echo.Context) error {
	var room Room
	roomID := c.Param("id")
	err := db_client.DB.QueryRow("SELECT rooms.name, users.name FROM rooms INNER JOIN users ON rooms.owner = users.id WHERE rooms.id = ?;", roomID).Scan(&room.Name, &room.Owner)
	if err != nil {
		return c.JSON(500, ErrMsg{
			Error: true,
			Msg:   "伺服器內部錯誤",
		})
	}

	return c.JSON(200, echo.Map{
		"ok":    true,
		"name":  room.Name,
		"owner": room.Owner,
	})
}
