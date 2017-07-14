package client

import (
	"github.com/garyburd/redigo/redis"
	"goJobSkills/constant"
)

var REDIS *redis.Pool

func Init() {
	REDIS =  &redis.Pool{
		MaxIdle: 		constant.MAX_IDLE,
		IdleTimeout: 	constant.IDLE_TIMEOUT,
		Dial: 		func () (redis.Conn, error) {
			return redis.Dial("tcp", constant.REDIS_SERVER)
		},
	}
}
