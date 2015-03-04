package logic

import (
	"time"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service"
)

func HandleBuyEnergy(gameSession *entity.GameSession, addEnergy int32, now time.Time) (messages pub.EntityMessagePairList) {
	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	userAddAttr.SetEnergy(userAddAttr.GetEnergy() + addEnergy)
	service.GetUserAddAttrService().Set(&userAddAttr)
	messages = net.EncodeBuyEnergyRecv(&userAddAttr)
	return
}
func HandleGotoDungeon(gameSession *entity.GameSession, subEnergy int32, now time.Time) (messages pub.EntityMessagePairList) {
	userAddAttr, ok := service.GetUserAddAttrService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	if userAddAttr.GetEnergy() < subEnergy {
		messages = net.EncodeGotoDungeonRecv(FAIL_STATUS, &userAddAttr)
		return
	}
	userAddAttr.SetEnergy(userAddAttr.GetEnergy() - subEnergy)
	service.GetUserAddAttrService().Set(&userAddAttr)
	messages = net.EncodeGotoDungeonRecv(SUCC_STATUS, &userAddAttr)
	return
}
