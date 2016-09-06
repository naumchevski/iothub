package main

import (
	"fmt"
	"log"

	"code.naumchevski.com/iot/command"
	"code.naumchevski.com/iot/connection"
	"code.naumchevski.com/iot/counter"
	"code.naumchevski.com/iot/thing"
)

const filepath = "./counters.json"
const publicFilepath = "./public-counters.json"

var privateCounters = counter.New(filepath)
var publicCounters = counter.New(publicFilepath)
var connections = connection.New()
var things = thing.New()

func handleUpdateCommand(s string) bool {
	cmd := command.Command(s)
	id, val, err := cmd.Handle()
	if err == nil {
		if counter.IsPrivateCounter(id) {
			if privateCounters.Update(id, val) {
				updatedConnections := connections.UpdateCounter(id, val)
				broadcastUpdatedCounters(updatedConnections)
			}
			return true
		} else {
			return false	
		}
	} else {
		log.Printf("cmd.Handle err: %s\n", err.Error())
		return false
	}
}

func broadcastUpdatedCounters(updatedConnections map[int]struct{}) {
	for id := range updatedConnections {
		t := things.Get(id)
		if t != nil {
			c := connections.Get(id)
			things.Broadcast(fmt.Sprintf("T%d", t.Id), t.BroadcastChan, c)
		}
	}
}
