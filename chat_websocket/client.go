package chat_websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var mu sync.Mutex

type Action int

const (
	ShowMembers Action = iota + 1
	Join
	SendMsg
	Leave
)

func (a Action) String() string {
	switch a {
	case ShowMembers:
		return "showMembers"
	case Join:
		return "join"
	case SendMsg:
		return "msg"
	case Leave:
		return "leave"
	default:
		return "unknown"
	}
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	name string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Msg struct {
	Action Action `json:"action"`
	Name   string `json:"name"`
	Msg    string `json:"msg"`
	RoomID int    `json:"roomID"`
}

func (c *Client) readPump() {
	defer func() {
		// leave the message to everyone still in the hub that one client has left
		memberLeave, err := json.Marshal(map[string]interface{}{
			"action": Leave,
			"name":   c.name,
		})
		if err != nil {
			log.Println(err)
		}
		c.hub.broadcast <- memberLeave

		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			// when connection can't read msg, break the loop and run the defer func() to close the connection and unregist the client
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// parse msg to struct and add some information in to it
		var parseMsg Msg
		json.Unmarshal(msg, &parseMsg)
		parseMsg.RoomID = c.hub.id
		parseMsg.Name = c.name
		msg, err = json.Marshal(parseMsg)
		// log.Printf("%#v", parseMsg)

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//hub close the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// Add queued chat messaages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1011, "websocket 升級失敗"))
		return err
	}

	// confirm which hub (hub id) will be connect, if hub_id not in the hubManager create hub
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1011, "字串轉數字失敗"))
		return err
	}

	// register a new hub and add the id into manager map
	mu.Lock()
	hub, ok := HM.hubs[roomID]
	if !ok {
		hub = NewHub(&HM, roomID)
		HM.hubs[roomID] = hub
		go hub.Run()
	}
	mu.Unlock()

	// increase the old member list in the hub and return to new member
	showOldMembers(conn, hub)

	// create client, allocate connection to the hub
	name := c.Get("userName")
	strName := fmt.Sprintf("%v", name)
	client := &Client{
		name: strName,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	// add new member name to everyone in the hub
	newMember, err := json.Marshal(map[string]interface{}{
		"action": Join,
		"name":   strName,
	})
	if err != nil {
		log.Println(err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1011, "轉換JSON失敗"))
		return err
	}
	client.hub.broadcast <- newMember

	go client.writePump()
	go client.readPump()
	return nil
}

func showOldMembers(conn *websocket.Conn, hub *Hub) {
	var memLst []string
	for client := range hub.clients {
		memLst = append(memLst, client.name)
	}
	conn.WriteJSON(map[string]interface{}{
		"action": ShowMembers,
		"data":   memLst,
	})
}
