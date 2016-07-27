// const.go

package ws

import (
	"time"
)

const (
	// 数据写入的最长时间
	writeWait = 10 * time.Second

	// 数据读入的最长时间
	pongWait = 60 * time.Second

	// ping的周期，一定要比pongWait小
	pingPeriod = (pongWait * 9) / 10

	// mssage的最大size
	maxMessageSize = 512
)
