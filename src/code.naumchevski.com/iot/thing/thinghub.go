package thing

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func New() *things {
	return &things{make(map[*websocket.Conn]*thing), sync.Mutex{}, 9}
}

func (t *things) Add(conn *websocket.Conn) *thing {
	(*t).mu.Lock()
	defer (*t).mu.Unlock()
	if conn != nil {
		broadcastChan := make(chan interface{}, broadcastChannelBuffer)
		done := make(chan struct{})
		
		// broadcast handler for voter websocket
		go func(conn *websocket.Conn, ch chan interface{}, done chan struct{}) {
			for {
				select {
				case value := <-ch:
					conn.SetWriteDeadline(time.Now().Add(broadcastWebSocketTimeout))
					if err := websocket.WriteJSON(conn, value); err != nil {
						log.Printf("error: send voter WS: %s", err)
					}
				case <-done:
					return
				}
			}
		}(conn, broadcastChan, done)
		
		(*t).counter++
		th := &thing{(*t).counter, broadcastChan, done}
		(*t).data[conn] = th
		return th
	}
	return nil
}

func (t *things) Get(id int) *thing {
	(*t).mu.Lock()
	defer (*t).mu.Unlock()
	for _ , t := range t.data {
		if t.Id == id {
			return t	
		}
	}
	return nil
}

func (t *things) Delete(conn *websocket.Conn) {
	(*t).mu.Lock()
	defer (*t).mu.Unlock()
	if conn != nil {
		delete((*t).data, conn)
	}
}

// broadcast to websocket client
func (t *things) Broadcast(client string, ch chan interface{}, value interface{}) {
	if cap(ch) != broadcastChannelBuffer {
		log.Printf("error: %s - use unbuffered broadcasting channel", client)
		return
	}
	// use channel only if the buffer is half full or half empty ;)
	if (cap(ch)+1)/2 > len(ch) {
		ch <- value
	} else {
		log.Printf("warn: %s - broadcasting on full channel\n", client)
	}
}
