// context.go

package ws

import (
	"encoding/json"
)

type Context struct {
	Conn *Connection

	Data *json.RawMessage

	Cmd string
}

func NewContext(conn *Connection, data *json.RawMessage) *Context {
	return &Context{conn, data, ""}
}

func (c *Context) Send(m []byte) {
	select {
	case c.Conn.send <- m:
	default:
		h.unregister <- c.Conn
	}
}

func (c *Context) EncodeWithCmd(b []byte) ([]byte, error) {
	m := Message{
		Cmd:  c.Cmd,
		Data: new(json.RawMessage),
	}
	*m.Data = b

	message, err := json.Marshal(m)
	return message, err
}

func (c *Context) String(message string) {
	c.Send([]byte(message))
}

func (c *Context) DecodeJson(v interface{}) error {
	err := json.Unmarshal(*c.Data, v)
	return err
}

func (c *Context) SetConnName(name string) {
	c.Conn.Name = name
	select {
	case h.setname <- c.Conn:
	default:
		h.unregister <- c.Conn
	}
}

func (c *Context) GetConnName() string {
	return c.Conn.Name
}
