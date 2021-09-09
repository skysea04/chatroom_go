package chat_websocket

type HubManager struct {
	hubs       map[int]*Hub
	register   chan int
	unregister chan *Hub
}

var HM HubManager = HubManager{
	hubs:       make(map[int]*Hub),
	register:   make(chan int),
	unregister: make(chan *Hub),
}

type Hub struct {
	manager    *HubManager
	id         int
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func CreateHubManager() *HubManager {
	return &HubManager{
		hubs:       make(map[int]*Hub),
		register:   make(chan int),
		unregister: make(chan *Hub),
	}
}

func NewHub(m *HubManager, id int) *Hub {
	return &Hub{
		manager:    m,
		id:         id,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *HubManager) Run() {
	for {
		select {
		case hub := <-m.unregister:
			if _, ok := m.hubs[hub.id]; ok {
				delete(m.hubs, hub.id)
				close(hub.broadcast)
				close(hub.register)
				close(hub.unregister)
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// log.Println("some one leave")
				delete(h.clients, client)
				close(client.send)
				if len(h.clients) == 0 {
					h.manager.unregister <- h
					return
				}
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
