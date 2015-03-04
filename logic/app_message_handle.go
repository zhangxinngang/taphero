package logic

import (
	"github.com/fanngyuan/link"
	"github.com/gogo/protobuf/proto"
	"time"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/resource"
)

func HandleMessage(msg []byte, session link.SessionAble, now time.Time) {
	absMsg, ok := pub.DecodeABSMessage(msg)
	if !ok {
		return
	}
	if len(absMsg.GetMsgList()) == 0 {
		return
	}

	log.Debugf("----->c2s[%s] msgid=%d,ip=%s",
		net.GetActionTypeName(pub.GetRealSendActionType(absMsg.GetMsgList()[0].GetActionType())),
		absMsg.GetMessageId(), session.Conn().RemoteAddr().String())
	// 发送消息回执
	if !SendAck(absMsg, session, now) {
		return
	}

	dispatch(session, absMsg, now)
	return
}

func dispatch(session link.SessionAble, c2sABSMsg pub.ABSMessage, now time.Time) {
	if len(c2sABSMsg.GetMsgList()) == 1 && pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()) == int32(pf.PFActionType_KeepSocketAuthSend) {
		submsg := pf.KeepSocketAuthSend{}
		proto.Unmarshal(c2sABSMsg.GetMsgList()[0].MessageBytes, &submsg)
		msgPair := HandleAuth(session, submsg, now)
		if len(msgPair) != 0 {
			gameSession, gameSessionok := resource.GetGameSessionByTcpSessionId(session.Id())
			if gameSessionok { // 把消息加入到 待检查队列中，
				gameSession.SendMsg(msgPair, now)
			}
		}
		return // auth 消息单独处理
	}

	gameSession, gameSessionok := resource.GetGameSessionByTcpSessionId(session.Id())
	if !gameSessionok { //
		log.Info("not login", session.Id(), session.Conn().RemoteAddr().String())
		return
	}
	msgHandle := entity.MsgHandle{
		// GameSession: &gameSession,
		// Session:           session,
		C2SABSMsg: c2sABSMsg,
		Now:       now,
		// HandleFunc: HandleUserMsgList,
		// RecvMsgHandleChan: make(chan entity.MsgHandle),
	}
	gameSession.PushMsgHandle(msgHandle)
	return
}

func HandleUserMsgList(msgHanle *entity.MsgHandle) {
	var s2cMsgPairList pub.EntityMessagePairList
	var isMessageReceiptSend bool = false
	c2sMsgId := msgHanle.C2SABSMsg.GetMessageId()
	for _, msgPair := range msgHanle.C2SABSMsg.GetMsgList() {
		tmpMsgList, tmpIsMessageReceiptSend := HandleUserMsg(c2sMsgId, msgHanle.GameSession, msgPair, msgHanle.Now)
		if tmpIsMessageReceiptSend {
			isMessageReceiptSend = true
		}

		if tmpMsgList != nil {
			s2cMsgPairList = append(s2cMsgPairList, tmpMsgList...)
		}
	}
	if !isMessageReceiptSend {
		msgHanle.GameSession.SetNextC2SMessageId(msgHanle.GameSession.GetNextC2SMessageId() + 1)
		log.Debugf("inscnextc2sMessage uin=%d,new c2smsgId=%d ip=%s", msgHanle.GameSession.Uin, msgHanle.GameSession.GetNextC2SMessageId(), msgHanle.GameSession.LinkSession.Conn().RemoteAddr().String())
	}

	msgHanle.GameSession.SendMsg(s2cMsgPairList, msgHanle.Now)

	return
}

func HandleUserMsg(c2sMsgId int32, gameSession *entity.GameSession, c2sMsgPair *pub.MessagePair, now time.Time) (messages pub.EntityMessagePairList, isMessageReceiptSend bool) {
	actionType := pub.GetRealSendActionType(c2sMsgPair.GetActionType())
	log.Infof("HandleUserMsg[%s]msgid=%d uin=%d,sessionid=%d,ip=%s",
		net.GetActionTypeName(actionType), c2sMsgId, gameSession.Uin, gameSession.LinkSession.Id(), gameSession.LinkSession.Conn().RemoteAddr().String())
	switch pf.PFActionType_PFActionTypeDetail(actionType) {
	case pf.PFActionType_MessageReceiptSend:
		isMessageReceiptSend = true
		submsg := pf.MessageReceiptSend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleMessageReceiptSend(submsg.GetReadMessageId(), gameSession, now)
	case pf.PFActionType_CheckVersion:
		submsg := pf.CheckVersionSend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleCheckVersion(gameSession, submsg.GetChannel(), submsg.GetVersion(), now)

	case pf.PFActionType_GetPayToken:
		messages = HandleGetPayToken(gameSession, now)
	case pf.PFActionType_OrderApply:
		submsg := pf.OrderApplySend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleOrderApply(gameSession, submsg.GetOrderList(), now)
	case pf.PFActionType_GameDataBase:
		messages = HandleGameDataBase(gameSession, now)
	case pf.PFActionType_SyncLogoutTime:
		submsg := pf.SyncLogoutTimeSend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleSyncLogoutTime(gameSession, submsg.GetStatus(), now)
		// case pf.PFActionType_SyncEnergy:
		// 	submsg := pf.SyncEnergySend{}
		// 	proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		// 	messages = HandleSyncEnergy(gameSession, submsg.GetEnergy(), submsg.GetEnergyAddTM(), now)
	case pf.PFActionType_GotoDungeon:
		submsg := pf.GotoDungeonSend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleGotoDungeon(gameSession, submsg.GetDungeonID(), now)
	case pf.PFActionType_BuyEnergy:
		submsg := pf.BuyEnergySend{}
		proto.Unmarshal(c2sMsgPair.MessageBytes, &submsg)
		messages = HandleBuyEnergy(gameSession, submsg.GetAddEnergy(), now)

	}
	return
}

// 发送消息回执
