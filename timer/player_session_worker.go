package timer

import (
	"github.com/0studio/scheduler"
	"time"
)

var PLAYER_SESSIOHN_TIME_WHEEL *scheduler.TimingWheel

const (
	CHECK_LOGIN_INFO_SECONDS = 40 * 60
)

func setUpPlayerSessionTimeWheel() {
	// 第二个参数至少要比defs.CHECK_LOGIN_INFO_SECONDS大 ,这里init 为 其10倍
	PLAYER_SESSIOHN_TIME_WHEEL = scheduler.NewTimingWheel(1*time.Second, CHECK_LOGIN_INFO_SECONDS*10)
}
func GetSessionTimeWheelChan(timeout time.Duration) <-chan struct{} {
	return PLAYER_SESSIOHN_TIME_WHEEL.After(timeout)
}

//
