package redisclient

import (
	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	Name string
	Ip   string
	Db   int
}

var G_RedisClient = make(map[string]redis.Conn)

func Redisinit(rcs []RedisConfig) error {
	for _, rc := range rcs {
		if conn, err := redis.Dial("tcp", rc.Ip, redis.DialDatabase(rc.Db)); err != nil {
			return err
		} else {
			G_RedisClient[rc.Name] = conn
		}
	}
	return nil
}

func RedisGetClient(name string) (conn redis.Conn, exist bool) {
	conn, exist = G_RedisClient[name]
	return
}
