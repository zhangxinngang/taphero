package logic

import (
	key "github.com/0studio/storage_key"
	"github.com/fanngyuan/link"
	"time"
	"zerogame.info/profile"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/resource"
	"zerogame.info/taphero/service"
)

func HandleAuth(session link.SessionAble, submsg pf.KeepSocketAuthSend, now time.Time) (messages pub.EntityMessagePairList) {
	var accountId, accountName string
	accountId = submsg.GetAccountID()
	accountName = submsg.GetAccountName()
	log.Infof("Auth accountid=%s,accountname=%s,channel=%d,uuid=%s", submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetChannelID(), submsg.GetUuid())
	if accountId == "" {
		log.Errorf("auth_accountid_name_empty [%s],[%s],[%s]", submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetUuid())
		accountId = submsg.GetUuid()
	}
	if accountName == "" {
		log.Errorf("auth_accountid_name_empty [%s],[%s],[%s]", submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetUuid())
		accountName = submsg.GetUuid()
	}

	if accountId == "" || accountName == "" {
		return
	}

	user, ok := service.GetUserService().Auth(submsg.GetAccountID(), submsg.GetAccountName(), conf.GetPlatform(), uint64(submsg.GetChannelID()), conf.GetServer())
	if !ok {
		user = newUser(
			service.GetUserService().GetNewUin(),
			submsg.GetAccountID(),
			submsg.GetAccountName(),
			submsg.GetUuid(),
			uint64(submsg.GetChannelID()),
			0, // gender
			submsg.GetOS(),
			submsg.GetOSVersion(),
			submsg.GetDeviceModel(),
			now)
		ok = service.GetUserService().Add(&user, now)
		log.Infof("create_user,result=%v,uin=%d,accountid=%s,accountname=%s,chan=%d,platform=%d", ok, user.GetUin(), submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetChannelID(), user.GetPlatform())
		if !ok {
			log.Error("create_user_err", user.GetUin(), submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetChannelID())
			log.Info("create_user_err", user.GetUin(), submsg.GetAccountID(), submsg.GetAccountName(), submsg.GetChannelID())
			return
		}
	}
	gameSession, gameSessionok := resource.GetGameSession(user.GetUin())
	log.Debugf("uin=%d,auth时客户端发过来的NextMessageId=%d(意思是说client 下一条发过来的消息其msgid为此值)", user.GetUin(), submsg.GetNextMessageId())
	if !gameSessionok {
		gameSession = entity.NewGameSession(user.GetUin(), session)
		gameSession.SetNextC2SMessageId(submsg.GetNextMessageId())
		resource.PutGameSession(user.GetUin(), gameSession)
		resource.PutUin(session.Id(), user.GetUin())
		log.Debugf("newgameSession,sessionid=%d,uin=%d,", session.Id(), gameSession.Uin)
		StartPlayerSessionScheduler(user.GetUin())
	} else {
		log.Debugf("auth时gameSession已存在,sessionid=%d,uin=%d,", session.Id(), gameSession.Uin)
		gameSession.Reconn(session, submsg.GetNextMessageId())
	}
	//处理漏单
	msg := handleUnhandledPayOrder(&user, now)
	messages = net.EncodeKeepSocketAuthRecv(gameSession.GetCurrentS2CMessageId())
	gameSession.SendMsg(msg, now)
	return

}
func newUser(uin uint64, accountId, accountName, uuid string, channel uint64, gender, os, osVersion int32, deviceModel string, now time.Time) (user profile.User) {
	user.SetUin(key.KeyUint64(uin))
	user.SetAccountId(accountId)
	user.SetAccountName(accountName)
	user.SetGender(gender)
	user.SetPlatform(conf.GetPlatform())
	user.SetChannel(channel)
	user.SetServer(conf.GetServer())
	user.SetUuid(uuid)
	user.SetOs(os)
	user.SetOsVersion(osVersion)
	user.SetDeviceModel(deviceModel)
	user.SetCreateTime(now)
	return
}
