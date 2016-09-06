package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.naumchevski.com/helper"
	"code.naumchevski.com/iot/counter"
	"code.naumchevski.com/iot/util"
)

const counterKeyPartLen = 4

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IoT Hub")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.RawQuery
	val := 0
	if id, ok := counter.GetPrivateCounterId(id, *publicCounters); ok {
		val = util.CounterToInt(privateCounters.Get(id))
	}
	fmt.Fprintf(w, strconv.Itoa(val))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if handleUpdateCommand(r.URL.RawQuery) {
		fmt.Fprintf(w, "success")
	} else {
		fmt.Fprintf(w, "error")
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// counter name
	c := fmt.Sprintf("%s-%s-%s-%s",
		helper.GenerateRandomStr(counterKeyPartLen),
		getNodeId(),
		helper.GenerateRandomStr(counterKeyPartLen),
		helper.GenerateRandomStr(counterKeyPartLen))
	privateCounters.Add(c, 0)

	// public counter name
	pc := fmt.Sprintf("%s-%s-%s",
		helper.GenerateRandomStr(counterKeyPartLen),
		getNodeId(),
		helper.GenerateRandomStr(counterKeyPartLen))
	publicCounters.Add(pc, c)

	type output struct {
		Public  string `json:"pub"`
		Private string `json:"priv"`
	}
	err := json.NewEncoder(w).Encode(&output{pc, c})
	if err != nil {
		log.Printf("cannot encode counters output: %s", err)
	}
}

func openHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	t := things.Add(ws)
	defer close(t.Done)
	if t == nil {
		return
	}
	defer things.Delete(ws)
	connections.Add(t.Id)
	defer connections.Delete(t.Id)
	params := r.URL.RawQuery
	for _, c := range strings.Split(params, ";") {
		if id, ok := counter.GetPrivateCounterId(c, *publicCounters); ok {
			val := privateCounters.Get(id)
			connections.SubscribeCounter(t.Id, id, util.CounterToInt(val))
		}
	}
	c := connections.Get(t.Id)
	things.Broadcast(fmt.Sprintf("T%d", t.Id), t.BroadcastChan, c)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if err != nil {
			log.Printf("client err: %s\n", err.Error())
			break
		} else {
			handleUpdateCommand(string(message))
		}
	}
}
