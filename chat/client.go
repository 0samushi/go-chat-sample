package chat

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id             int
	ws             *websocket.Conn
	removeClientCh chan *Client
	messageCh      chan []byte
}

func NewClient(ws *websocket.Conn, remove chan *Client, message chan []byte) *Client {
	return &Client{
		ws: ws,
		removeClientCh: remove,
		messageCh: message,
	}
}

func (c *Client) Start() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(msg)
		c.messageCh <- msg
	}
}

func (c *Client) Send(message []byte) {
	c.ws.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) Close() {
	c.ws.Close()
}
