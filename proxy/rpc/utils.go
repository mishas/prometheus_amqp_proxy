package rpc

import (
	"math/rand"
	"sync"
)

// chanMap is an atomic implementation of map[string]chan []byte.
type chanMap struct {
	m  map[string]chan []byte
	mu sync.RWMutex
}

func NewChanMap() *chanMap {
	return &chanMap{make(map[string]chan []byte), sync.RWMutex{}}
}

func (cm *chanMap) Add(k string, v chan []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.m[k] = v
}

func (cm *chanMap) Get(k string) (v chan []byte, ok bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	v, ok = cm.m[k]
	return
}

func (cm *chanMap) Delete(k string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.m, k)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}
