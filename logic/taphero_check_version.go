package logic

import (
	"time"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pub"
)

const (
	SUCC_STATUS = 0
	FAIL_STATUS = 1
)
const (
	CHANNEL_DEV = 112
	CHANNEL_IOS = 6
)

var DENY_CHANNEL = map[int32]bool{
// CHANNEL_DEV: true,
// CHANNEL_IOS: true,
}

// current client version 00100000
func HandleCheckVersion(gameSession *entity.GameSession, channel int32, version int32, now time.Time) (messages pub.EntityMessagePairList) {
	if DENY_CHANNEL[channel] {
		return net.EncodeCheckVersionRecv(FAIL_STATUS, FAIL_STATUS)
	}

	return net.EncodeCheckVersionRecv(SUCC_STATUS, SUCC_STATUS)
}
