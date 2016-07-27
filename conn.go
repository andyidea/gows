// conn.go

package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ws *websocket.Conn

	Name string

	send chan []byte
}

// 读数据
func (c *Connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		t, m, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Websocket error: %v", err)
			}
			break
		}

		// 不使用其它消息类型
		if t != websocket.TextMessage {
			log.Printf("Websocket error:connection massage type error.")
			break
		}

		var message Message
		err = json.Unmarshal(m, &message)
		if err != nil {
			log.Printf("Websocket error: message unmarshal error:%v\n %v", err, m)
			break
		}

		handlerFunc, ok := methods.Get(message.Cmd)
		if ok {
			context := NewContext(c, message.Data)
			context.Cmd = message.Cmd
			handlerFunc(context)
		} else {
			log.Println("Websocket error: not bind method named ", message.Cmd)
		}
	}
}

// 执行写入
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// 写数据
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		h.unregister <- c
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
