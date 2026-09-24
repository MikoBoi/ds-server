package websocket

import (
	"context"

	"github.com/coder/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c *Client) Send(ctx context.Context, data []byte) error {
	return c.conn.Write(ctx, websocket.MessageText, data)
}
