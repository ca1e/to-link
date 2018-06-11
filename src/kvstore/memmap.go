package kvstore

import (
	"errors"
)

var KVMap map[string]string = make(map[string]string)

type MemMap struct{
}

func (m *MemMap)Store(k, v string) error {
	if len(KVMap) > 1000000 {
		return errors.New("Limited")
	}
	KVMap[k] = v
	println("cur length:", len(KVMap))
	return nil
}

func (m *MemMap)Exist(k string) (string,bool){
	v, ok :=KVMap[k]
	return v, ok
}