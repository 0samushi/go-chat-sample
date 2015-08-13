package main

import (
//	"code.google.com/p/go.net/websocket"
	"net/http"
//	"fmt"
	"golang.org/x/net/websocket"
//	"image"
)

type Client struct {
	conn     *websocket.Conn
	incoming chan string
	outgoing chan string
}

func NewClient(ws *websocket.Conn) *Client {
	c := &Client{
		conn: ws,
		incoming: make(chan string),
		outgoing: make(chan string),
	}
	go c.Listen()

	return c
}

func (c *Client) Listen() {
	go c.read()
	go c.write()
}

func (c *Client) read() {
	for {
		var msg string
		websocket.Message.Receive(c.conn, &msg)
		c.outgoing <- msg
	}
}

func (c *Client) write() {
	for msg := range c.incoming{
		websocket.Message.Send(c.conn, msg)
	}
}

type Room struct {
	clients  map[*Client]bool
	incoming chan string
//	outgoing chan string
	joins    chan *websocket.Conn
}

func NewRoom() *Room {
	r := &Room{
		clients:make(map[*Client]bool),
		incoming:make(chan string),
//		outgoing:make(chan string),
		joins: make(chan *websocket.Conn),
	}

	go r.listen()
	return r
}

func (r *Room) listen() {
	for {
		select {
		case conn := <-r.joins:
			r.join(conn)
		case msg := <-r.incoming:
			r.broadcast(msg)
		}
	}
}

func (r *Room) broadcast(msg string) {
	for client, _ := range r.clients {
		client.incoming <- msg
	}
}

func (r *Room) join(conn *websocket.Conn) {
	c := NewClient(conn)
	r.clients[c] = true
	go func() {
		for {
			r.incoming <- <-c.outgoing
		}}()
}

func echoHandler(ws *websocket.Conn) {

////	for {
		room.joins <- ws
//	}
}

var room *Room

func main() {
	room = NewRoom()
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{Handler: websocket.Handler(echoHandler)}
		s.ServeHTTP(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
