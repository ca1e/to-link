package kvstore

import (
	"github.com/gomodule/redigo/redis"
)

type Redigo struct{
}

//https://studygolang.com/articles/3029
var (
	MAX_POOL_SIZE = 20
	redisPoll chan redis.Conn 
)

func putRedis(conn redis.Conn) {
    if redisPoll == nil {
        redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
    }
    if len(redisPoll) >= MAX_POOL_SIZE {
        conn.Close()
        return
    }
    redisPoll <- conn
} 

func initRedis(network, address string) redis.Conn {
    if len(redisPoll) == 0 {
        redisPoll = make(chan redis.Conn, MAX_POOL_SIZE)
        go func() {
            for i := 0; i < MAX_POOL_SIZE/2; i++ {
                c, err := redis.Dial(network, address)
                if err != nil {
                    panic(err)
                }
                putRedis(c)
            }
        } ()
    }
    return <-redisPoll
}

func (r *Redigo)Store(k, v string) error {
	c := initRedis("tcp", "127.0.0.1:6379")
	vr, err := c.Do("set", k, v)
	println("set redis:", vr)
	if err != nil {
	  return err
    }
    c.Do("EXPIRE", k, 60*60*24*30)
	leng, err := redis.Int64(c.Do("dbsize"))
	if err != nil {
	  return err
	}
	println("cur length:", leng)
	return nil
}

func (r *Redigo)Exist(k string) (string,bool){
	c := initRedis("tcp", "127.0.0.1:6379")
	ok, _ := redis.Int64(c.Do("exists", k))
	v, _ := redis.String(c.Do("get", k))
	return v, (ok==1)
}
