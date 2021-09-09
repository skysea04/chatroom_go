package main

import (
	"html/template"
	"io"
	"main/chat_websocket"
	"main/controllers"
	"main/db_client"
	"main/views"

	"github.com/labstack/echo/v4"
)

var manager = chat_websocket.CreateHubManager()

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	db_client.InitialiseDBConnection()
	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.GET("/signup", views.Signup)
	e.POST("/user", controllers.PostUser)
	e.POST("/user-entry", controllers.LoginUser)

	apis := e.Group("/api")
	apis.Use(controllers.JwtGateKeeper)
	apis.POST("/room", controllers.PostRoom)
	apis.GET("/rooms", controllers.GetRooms)
	apis.GET("/my/rooms", controllers.GetMyRooms)
	apis.GET("/room/:id", controllers.GetRoomInfo)

	// view 瀏覽畫面
	v := e.Group("")
	v.Use(controllers.JwtGateKeeper)
	v.GET("/", views.Index)
	v.GET("/create-room", views.CreateRoom)
	v.GET("/my/rooms", views.MyRooms)
	v.GET("/chatroom/:id", views.ChatRoom)

	// websocket
	// manager := chat_websocket.CreateHubManager()
	go chat_websocket.HM.Run()
	v.GET("/ws/:id", chat_websocket.ServeWs)

	e.Logger.Fatal(e.Start(":8000"))
}
