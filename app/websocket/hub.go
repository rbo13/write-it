package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Hub is our handler
// that represents a Hub
type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
}

// NewHub is our constructor that
// returns an instance of Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

}

// Run runs the websocket server
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

var upgrader = websocket.Upgrader{
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebsocket handles websocket connection
func (hub *Hub) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := NewClient(hub, socket)
	hub.clients = append(hub.clients, client)
	hub.register <- client
	client.run()
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())
	// TODO:: implement properly onConnect
	// Make list of all users
	// users := []message.User{}
	// for _, c := range hub.clients {
	// 	users = append(users, message.User{ID: c.id, Color: c.color})
	// }
	// Notify user joined
	// hub.send(message.NewConnected(client.color, users), client)
	// hub.broadcast(message.NewUserJoined(client.id, client.color), client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())
	client.close()
	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}
	// Delete client from list
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]
	// Notify user left
	// hub.broadcast(message.NewUserLeft(client.id), nil)
}

func (hub *Hub) onMessage(data []byte, client *Client) {
	log.Println("Implement onMessage")
	// kind := gjson.GetBytes(data, "kind").Int()
	// if kind == message.KindStroke {
	// 	var msg message.Stroke
	// 	if json.Unmarshal(data, &msg) != nil {
	// 		return
	// 	}
	// 	msg.UserID = client.id
	// 	hub.broadcast(msg, client)
	// } else if kind == message.KindClear {
	// 	var msg message.Clear
	// 	if json.Unmarshal(data, &msg) != nil {
	// 		return
	// 	}
	// 	msg.UserID = client.id
	// 	hub.broadcast(msg, client)
	// }
}
