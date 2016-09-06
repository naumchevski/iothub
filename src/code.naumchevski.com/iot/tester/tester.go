package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

var origin = flag.String("origin", "http://localhost:80/", `address of server to connect`)
var open = flag.String("open", "", `list of counter splitted by ;. e.g. -open="ws://localhost:8080/open?xxxx-yyyy-xxxx-xxxx"`)
var set = flag.String("set", "", `set counter value e.g. -set="xxxx-yyyy-xxxx-xxxx=10"`)

func main() {
	flag.Parse()

	origin := *origin
	url := *open
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if len(*set) > 0 {
			if _, err := ws.Write([]byte(*set)); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Send: %s\n", *set)
		}
	}()

	if strings.Contains(*open, "?") {
		var msg = make([]byte, 1024)
		var n int
		for {
			if n, err = ws.Read(msg); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Received: %s\n", msg[:n])
			}
		}
	}
	
	time.Sleep(1 * time.Second)
}
