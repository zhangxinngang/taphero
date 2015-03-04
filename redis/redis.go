package redis

import (
	"fmt"
	redisapi "github.com/0studio/redisapi"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/log"
)

func SetUpRedis(addr string, maxActive, maxIdle int, wait bool) (client redisapi.Redis) {
	var err error
	client, err = redisapi.InitRedisClient(addr, maxActive, maxIdle, wait)
	if err != nil {
		log.Error("setup error:", addr, err)
	}
	log.Infof("waiting ping  %s result...", addr)
	if !conf.IsModeTest() {
		fmt.Printf("waiting ping  %s result...", addr)
	}
	if client.Ping() {
		if !conf.IsModeTest() {
			log.Info(" succ!!!")
			fmt.Println(" succ!!!")
		}
	} else {
		if !conf.IsModeTest() {
			fmt.Println("failed ")
			log.Error("failed")

			panic("redis_config")

		}
	}
	return
}
