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
		if conn, err := redis.Dial("tcp", rc.ip, redis.DialDatabase(rc.db)); err != nil {
			return err
		} else {
			G_RedisClient[rc.name] = conn
		}
	}
	return nil
}

func RedisGetClient(name string) (conn redis.Conn, exist bool) {
	conn, exist = G_RedisClient[name]
	return
}
