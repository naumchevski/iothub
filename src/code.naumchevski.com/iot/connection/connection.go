package connection

import (
	"sync"
)

type connection map[string]int

type connections struct {
	data map[int]*connection
	mu   sync.Mutex
}
