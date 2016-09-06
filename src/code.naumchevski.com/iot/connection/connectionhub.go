package connection

import (
	"sync"
)

func New() *connections {
	return &connections{make(map[int]*connection), sync.Mutex{}}
}

func (cs *connections) Add(id int) {
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	(*cs).data[id] = &connection{}
}

func (cs *connections) Delete(id int) {
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	delete((*cs).data, id)
}

func (cs *connections) Get(id int) connection {
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	return *(*cs).get(id)
}

func (cs *connections) SubscribeCounter(id int, counter string, value int) {
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	c := cs.get(id)
	if c != nil {
		if c == nil {
			n := make(connection)
			c = &n
		}
		(*c)[counter] = value
	}
}

func (cs *connections) UnsubscribeCounter(id int, counter string) {
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	c := cs.get(id)
	if c != nil {
		delete(*c, counter)
	}
}

func (cs *connections) UpdateCounter(counter string, value int) map[int]struct{} {
	var updatedConnections = make(map[int]struct{})
	(*cs).mu.Lock()
	defer (*cs).mu.Unlock()
	for c := range (*cs).data {
		if cs.update(c, counter, value) {
			updatedConnections[c] = struct{}{}
		}
	}
	return updatedConnections
}

func (cs *connections) update(id int, counter string, value int) bool {
	c := cs.get(id)
	if c != nil {
		if _, ok := (*c)[counter]; ok {
			(*c)[counter] = value
			return true
		}
	}
	return false
}

func (cs *connections) get(id int) *connection {
	return (*cs).data[id]
}
