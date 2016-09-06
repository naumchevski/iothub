package thing

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	broadcastChannelBuffer    int           = 100 // biffer size
	broadcastWebSocketTimeout time.Duration = 1 * time.Second   // timeout in seconds
)

type thing struct {
	Id            int
	BroadcastChan chan interface{}
	Done chan struct{}
}

type things struct {
	data    map[*websocket.Conn]*thing
	mu      sync.Mutex
	counter int
}
