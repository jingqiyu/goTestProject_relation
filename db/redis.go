package db

import (
	"github.com/gomodule/redigo/redis"
	"relation/conf"
	"time"
)

var (
	redisPool *redis.Pool
)

func InitRedisPool(configType conf.ConfigType) {
	var (
		maxIdle int
		addr string
		idleTimeOut int64
	)
	maxIdle = configType.DB.Redis.MaxIdle
	idleTimeOut = int64(configType.DB.Redis.Timeout)
	addr = "127.0.0.1:6379"
	redisPool = &redis.Pool{
		MaxIdle:maxIdle,
		IdleTimeout: time.Duration(idleTimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",addr)
		},
	}
}

func GetRedisClient() redis.Conn{
	return redisPool.Get()
}