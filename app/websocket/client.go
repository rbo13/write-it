package websocket

import (
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/rbo13/write-it/app/generate"
)

// Client is our handler
// for client socket connection
type Client struct {
	id       string
	hub      *Hub
	color    string
	socket   *websocket.Conn
	outbound chan []byte
}

// Create a Version 4 UUID, panicking on error.
// Use this form to initialize package-level variables.
var u1 = uuid.Must(uuid.NewV4())

// NewClient is our constructor
// that returns an instance of Client
func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	uuID, _ := uuid.NewV4()
	uuIDStr := uuID.String()
	return &Client{
		id:       uuIDStr,
		color:    generate.Color(),
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()
	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
		client.hub.onMessage(data, client)
	}
}

func (client *Client) write() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (client Client) run() {
	go client.read()
	go client.write()
}

func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}
