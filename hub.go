// hub.go

package ws

import (
	"encoding/json"
	"log"
	"runtime"
	"time"
)

type hub struct {
	// 记录所有的连接
	connections map[*Connection]string

	// 广播
	broadcast chan []byte

	// 根据name发送信息
	sendchan chan *NameData

	//设置名称
	setname chan *Connection

	// 注册连接的请求
	register chan *Connection

	// 注销连接的请求
	unregister chan *Connection
}

var h = hub{

	broadcast:   make(chan []byte),
	sendchan:    make(chan *NameData),
	setname:     make(chan *Connection),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]string),
}

func (h *hub) Run() {
	go func() {
		timer := time.Tick(time.Second * 30)
		for {
			select {
			case c := <-h.register:
				h.onRegister(c)
			case c := <-h.unregister:
				h.onUnRegister(c)
			case m := <-h.broadcast:
				h.onBroadcast(m)
			case n := <-h.sendchan:
				h.onSendByName(n)
			case c := <-h.setname:
				h.onSetName(c)
			case <-timer:
				h.onConnectionInfo()
			}
		}
	}()

}

// 注册链接
func (h *hub) onRegister(c *Connection) {
	h.connections[c] = ""
}

// 注销连接
func (h *hub) onUnRegister(c *Connection) {
	if _, ok := h.connections[c]; ok {
		delete(h.connections, c)
		close(c.send)
	}
}

// 广播
func (h *hub) onBroadcast(m []byte) {
	for c, name := range h.connections {
		if name != "" {
			select {
			case c.send <- m:
			default:
				close(c.send)
				delete(h.connections, c)
			}
		}
	}
}

// 根据name发送信息
func (h *hub) onSendByName(n *NameData) {
	for c, name := range h.connections {
		if n.Name == name {
			select {
			case c.send <- n.Data:
			default:
				close(c.send)
				delete(h.connections, c)
			}
		}
	}
}

func (h *hub) onSetName(c *Connection) {
	if _, ok := h.connections[c]; ok {
		h.connections[c] = c.Name
	}
}

// 打印连接信息
func (h *hub) onConnectionInfo() {
	numGoroutine := runtime.NumGoroutine()
	numConnections := len(h.connections)
	log.Printf("Goroutine %v Connection %v", numGoroutine, numConnections)
}

// 绑定
func (h *hub) Bind(r string, hf HandlerFunc) {
	methods = append(methods, Method{r, hf})
}

func (h *hub) SendMsgByConnName(connName string, msg []byte) {
	nameData := &NameData{connName, msg}
	h.sendchan <- nameData

}

func (h *hub) SendDataByConnName(connName string, cmd string, data []byte) {
	mssage := &Message{cmd, new(json.RawMessage)}
	*mssage.Data = data
	jresult, err := json.Marshal(mssage)
	if err != nil {
		log.Println("SendDataByConnName Json Marshal failed :", err.Error())
	}

	h.SendMsgByConnName(connName, jresult)
}

func (h *hub) Broadcast(cmd string, data []byte) {
	mssage := &Message{cmd, new(json.RawMessage)}
	*mssage.Data = data
	jresult, err := json.Marshal(mssage)
	log.Println(jresult)
	if err != nil {
		log.Println("Broadcast Json Marshal failed :", err.Error())
	}
	h.broadcast <- jresult
}
