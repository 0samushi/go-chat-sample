package chat

import "golang.org/x/net/websocket"

type Client struct {
	Id             int
	ws             *websocket.Conn
	removeClientCh chan *Client
	messageCh      chan string
}

func NewClient(ws *websocket.Conn, remove chan *Client, message chan string) *Client {
	return &Client{
		ws: ws,
		removeClientCh: remove,
		messageCh: message,
	}
}

func (c *Client) Start() {
	for {
		var msg string
		websocket.Message.Receive(c.ws, &msg)
		c.messageCh <- msg
	}
}

func (c *Client) Send(message string) {
	websocket.Message.Send(c.ws, message)
}

func (c *Client) Close() {
	c.ws.Close()
}
