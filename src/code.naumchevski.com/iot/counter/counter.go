package counter

import (
	"sync"
)

type counters struct {
	data map[string]interface{}
	mu   sync.Mutex
	autoSave bool
}
