package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"fmt"
)


type Server struct {
	clientCount    int
	clients        map[int]*Client
	addClientCh    chan *Client
	removeClientCh chan *Client
	messageCh      chan []byte
}

func NewServer() *Server {
	return &Server{
		clientCount:0,
		clients:map[int]*Client{},
		addClientCh:make(chan *Client),
		removeClientCh:make(chan *Client),
		messageCh: make(chan []byte),
	}
}

func (server *Server) addClient(client *Client) {
	server.clientCount++
	client.Id = server.clientCount
	server.clients[client.Id] = client
}

func (server *Server) removeClient(client *Client) {
	delete(server.clients, client.Id)
}

func (server *Server) sendMessage(message []byte) {
	for _, client := range server.clients {
		c := client
		go func() {c.Send(message)}()
	}
}

func (server *Server) Start() {
	for {
		select {
		case client := <-server.addClientCh:
			server.addClient(client)
		case client := <-server.removeClientCh:
			server.removeClient(client)
		case message := <-server.messageCh:
			server.sendMessage(message)
		}
	}
}


var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}


func (server *Server) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("create server!!!")
	client := NewClient(ws, server.removeClientCh, server.messageCh)
	server.addClientCh <- client
	client.Start()
}



