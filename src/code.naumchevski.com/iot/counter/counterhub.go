package counter

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

func New(filepath string) *counters {
	c := &counters{make(map[string]interface{}), sync.Mutex{}, false}
	go c.autoHandler(filepath)
	return c
}

func (c *counters) Add(id string, defaultVal interface{}) {
	c.mu.Lock()
	c.data[id] = defaultVal
	c.autoSave = true
	c.mu.Unlock()
}

func (c counters) Get(id string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[id]; ok {
		return c.data[id]
	}
	
	return nil
}

func (c *counters) Update(id string, val interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[id]; ok {
		if c.data[id] != val {
			c.data[id] = val
			c.autoSave = true
			return true
		}
	}
	return false

}

func (c *counters) autoHandler(path string) {
	// load
	if err := c.load(path); err != nil {
		log.Printf("error: load counters - %s", err.Error())
	}

	// auto save
	for {
		select {
		case <-time.After(10 * time.Second):
			c.mu.Lock()
			if c.autoSave {
				if err := c.save(path); err != nil {
					log.Printf("error: save counters - %s", err.Error())
				}
				c.autoSave = false
			}
			c.mu.Unlock()
		}
	}
}

func (c *counters) save(filename string) error {
	var err error
	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		err = json.NewEncoder(fl).Encode(&c.data)
		fl.Close()
	}
	return err
}

// load answer
func (c *counters) load(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	fl, err := os.Open(filename)
	if err == nil {
		tmp := make(map[string]interface{})
		err = json.NewDecoder(fl).Decode(&tmp)
		fl.Close()
		if tmp != nil {
			c.data = tmp
		}
	}
	return err
}
