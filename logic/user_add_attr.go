package logic

import (
	"time"
	"zerogame.info/taphero/defs"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service"
)

func HandleGameDataBase(gameSession *entity.GameSession, now time.Time) (messages pub.EntityMessagePairList) {
	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	messages = net.EncodeGameDataBaseRecv(userAddAttr)
	return
}

const (
	STATUS_SYNC_LOGOUT_TIME_CLEAR = 0
	STATUS_SYNC_LOGOUT_TIME_SET   = 1
)

func HandleSyncLogoutTime(gameSession *entity.GameSession, status int32, now time.Time) (messages pub.EntityMessagePairList) {
	//   optional int32 status = 1;		// 0:时间清零 1:设置服务器当前时间
	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	if STATUS_SYNC_LOGOUT_TIME_CLEAR == status {
		userAddAttr.SetLastOffTime(defs.UninitedTime)
		service.GetUserAddAttrService().Set(&userAddAttr)
	} else if STATUS_SYNC_LOGOUT_TIME_SET == status {
		userAddAttr.SetLastOffTime(now)
		service.GetUserAddAttrService().Set(&userAddAttr)
	}
	return
}

// func HandleSyncEnergy(gameSession *entity.GameSession, energy, energyTime int32, now time.Time) (messages pub.EntityMessagePairList) {
// 	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
// 	if !ok {
// 		return
// 	}
// 	userAddAttr.SetEnergy(energy)
// 	userAddAttr.SetEnergyTime(time.Unix(int64(energyTime), 0))
// 	service.GetUserAddAttrService().Set(&userAddAttr)
// 	return
// }
func onTcpClosed4LastOffLineTime(gameSession *entity.GameSession, closedTime time.Time, now time.Time) {
	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	if userAddAttr.GetLastOffTime() == defs.UninitedTime {
		userAddAttr.SetLastOffTime(closedTime)
		service.GetUserAddAttrService().Set(&userAddAttr)
	}
}
