package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/websocket"
)

var port = flag.String("port", ":80", "http server port")
var nodeId = flag.String("id", "iot0", "node id")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()

	http.HandleFunc("/open", openHandler)

	http.HandleFunc("/set", setHandler)

	http.HandleFunc("/get", getHandler)
	
	http.HandleFunc("/create", createHandler)

	http.HandleFunc("/", homeHandler)

	// demo
	startDemo()

	err := http.ListenAndServe(*port, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func getNodeId() string {
	return *nodeId;
}
