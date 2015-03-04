package redis_msg

// import (
// 	"zerogame.info/taphero/defs"
// 	"zerogame.info/taphero/redis"
// )

// func PubCenterMsg(value redis.RedisQueueValue) {
// 	// 往中心服的redis 队列时扔数据
// 	// 前提是msgRedisConn 初始化了
// 	// 否则不要调此函数
// 	redis.PubMsg(msgRedisConn, defs.REDIS_QUEUE_NAME_CENTER_SERVER, value)
// }

// func SendServerMsg(serverList entity.ServerList, value redis.RedisQueueValue) {
// 	for idx, _ := range serverList {
// 		SendServerProcessMsg(uint64(serverList[idx].Index), uint64(serverList[idx].ProcessIndex), value)
// 	}
// }

// func SendServerProcessMsg(server uint64, process uint64, value redis.RedisQueueValue) {
// 	queueName := fmt.Sprintf("notify_queue_%d_%d_%d", conf.GetDefaultPlatformId(), server, process)
// 	redis.SendMsg(msgRedisConn, queueName, value)
// }
