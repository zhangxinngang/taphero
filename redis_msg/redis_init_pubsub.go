package redis_msg

import (
	"github.com/0studio/redisapi"
	"github.com/garyburd/redigo/redis"
	"zerogame.info/taphero/conf"
)

func createPubSubConnPool() *redis.Pool {
	addr, maxActive, maxIdle := conf.GetServerRedisConfig()
	return redisapi.CreateRedisPool(addr, maxActive, maxIdle, true)
}
