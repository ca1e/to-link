package kvstore

import (
	"errors"
	"fmt"
	"sync"
)

var KVMap map[string]string = make(map[string]string)

type MemMap struct {
	lock sync.RWMutex
}

func (m *MemMap) Store(k, v string) error {
	if len(KVMap) > 1000000 {
		return errors.New("Limited")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	KVMap[k] = v
	fmt.Println("cur length:", len(KVMap))
	return nil
}

func (m *MemMap) Exist(k string) (string, bool) {
	v, ok := KVMap[k]
	return v, ok
}
