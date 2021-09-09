package kvstore

import (
	"errors"
	"sync"
)

var KVMap map[string]string = make(map[string]string)
var lock sync.RWMutex

type MemMap struct{
}

func (m *MemMap)Store(k, v string) error {
	if len(KVMap) > 1000000 {
		return errors.New("Limited")
	}
	
	lock.Lock()
	defer lock.Unlock()
	
	KVMap[k] = v
	println("cur length:", len(KVMap))
	return nil
}

func (m *MemMap)Exist(k string) (string,bool){
	v, ok :=KVMap[k]
	return v, ok
}
