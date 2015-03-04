package redis_msg

//同一服务器id的 都会收到此消息 ，通过redis pub sub 发布定况消息

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/redisapi"
	key "github.com/0studio/storage_key"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/defs"
	"zerogame.info/taphero/redis"
	"zerogame.info/taphero/resource"
)

var (
	redisSubChanChecker redisapi.RedisSubChanChecker
)

func setUpServerMsgChecker() {
	queueName := fmt.Sprintf("tap_hero_%d_%d", conf.GetPlatform(), conf.GetServer())
	redisSubChanChecker = redisapi.NewRedisSubChanChecker(createPubSubConnPool, queueName, serverMsgCheck)

	redisSubChanChecker.Start()
}
func stopServerMsgChecker() {
	go redisSubChanChecker.Stop()
}

func serverMsgCheck(dataFromRedis interface{}) {
	var value redis.RedisQueueValue
	err := json.Unmarshal(dataFromRedis.([]byte), &value)
	if err != nil {
		fmt.Println("redisQueue read err,", err)
		return
	}
	if value.Type == defs.NOTIFY_TYPE_PAY_LEAK_ORDER { // 漏单处理
		uin := value.Uin
		var keyUin key.KeyUint64
		if !keyUin.FromString(uin) {
			return
		}
		gameSession, ok := resource.GetGameSession(keyUin)
		if !ok {
			return
		}
		gameSession.TryCheckPayOrder()

	}

	//
}
