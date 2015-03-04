package app

import (
	"fmt"
	"github.com/fanngyuan/link"
	"time"
	"zerogame.info/taphero/conf"
	// "zerogame.info/taphero/design"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/logic"
	"zerogame.info/taphero/redis_msg"
	"zerogame.info/taphero/resource"
	"zerogame.info/taphero/service"
	"zerogame.info/taphero/timer"
)

func InitResource(mode string, server, process uint64, locale string) {
	conf.Init(mode, server, process, locale)
	log.Init()
	timer.Init()
	// design.SetUp()
	service.Init()
	resource.Init()
	redis_msg.SetUp()
	fmt.Println("setup done!!!")
	// timer.SetUpTcpCloserTimeWheel()
}
func CloseTcp(session link.SessionAble, closedTime time.Time, reason interface{}) {
	session.Conn().Close()
	session.Close(reason)
	logic.OnTcpClosed(session, closedTime, reason)
}
