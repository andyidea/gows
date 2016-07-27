// protocol.go

package ws

import (
	"encoding/json"
)

type Message struct {
	Cmd  string           `json:"cmd"`
	Data *json.RawMessage `json:"data"` // must be *json.RawMessage
}

type NameData struct {
	Name string
	Data []byte
}
