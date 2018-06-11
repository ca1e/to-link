package kvstore

type kkvv interface {
	Store(k, v string) error
	Exist(k string) (string,bool)
}

var kv kkvv

func init() {
	kv = &Redigo{}
}

func Store(k, v string) error {
	return kv.Store(k, v)
}

func Exist(k string) (string,bool){
	return kv.Exist(k)
}