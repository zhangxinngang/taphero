package redis_msg

import (
	"github.com/0studio/redisapi"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/redis"
)

var msgRedisConn redisapi.Redis

func SetUp() {
	setUpMsgRedis()
	// setUpMsgPubSubRedis()

	// setUpRedisQueueChecker()
	setUpServerMsgChecker()
}

func setUpServer() {
}

func Stop() {
	log.Info("shutdowning before  gmtool.DoShutDown() ...")
	// stopRedisQueueChecker()
	stopServerMsgChecker()
	log.Info("shutdowning after  gmtool.DoShutDown() ...")

}

func setUpMsgRedis() {
	addr, maxActive, maxIdle := conf.GetServerRedisConfig()
	msgRedisConn = redis.SetUpRedis(addr, maxActive, maxIdle, true)
}

func PushMsg(key string, value interface{}) error {
	return msgRedisConn.Lpush(key, value)
}
