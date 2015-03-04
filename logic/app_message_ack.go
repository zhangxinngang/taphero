package logic

import (
	"github.com/fanngyuan/link"
	"time"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/resource"
)

func SendAck(c2sABSMsg pub.ABSMessage, session link.SessionAble, now time.Time) bool {
	// log.Debugf("尝试发送消息回执给客户端%s,ip=%s",
	// 	net.GetActionTypeName(pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType())),
	// 	session.Conn().RemoteAddr().String())
	// return 此函数后，还需不要往下处理客户端发过来的消息
	if len(c2sABSMsg.GetMsgList()) == 0 {
		return false
	}
	if pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()) == int32(pf.PFActionType_MessageReceiptSend) {
		// 如果收到是客户端的回执消息， 则不再发回执， 否则无限循环了
		return true
	}
	if pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()) == int32(pf.PFActionType_KeepSocketAuthSend) {
		// 如果是认证消息 一定发回执
		log.Debugf("<----s2c[%s回执],%s c2sMsgId=%d,s2cMsgId=%d,ip=%s",
			net.GetActionTypeName(int32(pf.PFActionType_MessageReceiptRecv)),
			net.GetActionTypeName(pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType())),
			c2sABSMsg.GetMessageId(), 0,
			session.Conn().RemoteAddr().String())

		data := net.EncodeABSMessageWithHeader(pub.NewMsg(0, net.EncodeMessageReceiptRecv(c2sABSMsg.GetMessageId()), "", 0, now), now) //
		session.SendPacket(net.WritePacket(data))
		return true
	}

	gameSession, gameSessionok := resource.GetGameSessionByTcpSessionId(session.Id())
	if !gameSessionok { //玩家未登录， 且此条消息不是认证消息（上面已经判断过,直接返回，此种消息不处理）
		log.Info("not_login_and_not_auth_send_msg_do_not_send_ack ", session.Id(), pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()), session.Conn().RemoteAddr().String(), session.Conn().RemoteAddr().String())
		log.Error("not_login_and_not_auth_send_msg_do_not_send_ack", session.Id(), pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()), session.Conn().RemoteAddr().String(), session.Conn().RemoteAddr().String())
		return false
	}
	if c2sABSMsg.GetMessageId() < gameSession.GetNextC2SMessageId() {
		// 如果客户端发过来的消息id比 客户端应该发过来的消息id小，则只给其发回执， 不处理此消息逻辑
		log.Debug("客户端发过来的messageId比预期小,我只发回执 ，不予处理",
			net.GetActionTypeName(pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType())),
			session.Conn().RemoteAddr().String())
		log.Error("客户端发过来的messageId比预期小,我只发回执 ，不予处理",
			net.GetActionTypeName(pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType())),
			gameSession.Uin,
			session.Conn().RemoteAddr().String())
		data := net.EncodeABSMessageWithHeader(pub.NewMsg(0, net.EncodeMessageReceiptRecv(c2sABSMsg.GetMessageId()), "", 0, now), now) //
		session.SendPacket(net.WritePacket(data))
		return false
	}
	if c2sABSMsg.GetMessageId() == gameSession.GetNextC2SMessageId() {
		// 发回执， 处理消息
		log.Debugf("<----s2c[%s回执],%s MsgId=%d,ip=%s",
			net.GetActionTypeName(pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType())),
			net.GetActionTypeName(int32(pf.PFActionType_MessageReceiptRecv)),
			c2sABSMsg.GetMessageId(),
			session.Conn().RemoteAddr().String())
		data := net.EncodeABSMessageWithHeader(pub.NewMsg(0, net.EncodeMessageReceiptRecv(c2sABSMsg.GetMessageId()), "", 0, now), now) //
		// gameSession.InscNextC2SMessageId()
		session.SendPacket(net.WritePacket(data))
		return true
	}
	if c2sABSMsg.GetMessageId() > gameSession.GetNextC2SMessageId() {
		// 不发回执， 也不处理消息
		log.Debugf("客户端发过来的messageId比预期的大，不予回执 延后处理  %sactiontyp=%d,client c2smmsgId=%d,server c2sMsgId=%d",
			session.Conn().RemoteAddr().String(),
			pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()), c2sABSMsg.GetMessageId(), gameSession.GetNextC2SMessageId())

		// log.Warnf("客户端发过来的messageId比预期的大，不予回执 延后处理  %sactiontyp=%d,client c2smmsgId=%d,server c2sMsgId=%d",
		// session.Conn().RemoteAddr().String(),
		// pub.GetRealSendActionType(c2sABSMsg.GetMsgList()[0].GetActionType()), c2sABSMsg.GetMessageId(), gameSession.GetNextC2SMessageId())

		msgHandle := entity.MsgHandle{
			// GameSession: &gameSession,
			C2SABSMsg: c2sABSMsg,
			Now:       now,
		}

		gameSession.PushAckMsgHandle(msgHandle)
		return false
	}

	return false
}

func HandleMessageReceiptSend(msgId int32, gameSession *entity.GameSession, now time.Time) (messages pub.EntityMessagePairList) {
	// var hasMatchMsgId bool
	// 客户端收到我发给它的消息，可以将此消息从s2c 消息队列里去除了
	log.Debugf("客户端发过来已收到我服务器消息的回执消息 uin=%d,msgid=%d\n", gameSession.Uin, msgId)
	// if msgId != gameSession.GetCurrentS2CMessageId() {
	// 	fmt.Printf("HandleMessageReceiptSend s2cMsgId not match uin=%d,clientId=%d,serverid=%d\n", gameSession.Uin, msgId, gameSession.GetCurrentS2CMessageId())
	// 	log.Errorf("HandleMessageReceiptSend s2cMsgId not match uin=%d,clientId=%d,serverid=%d\n", gameSession.Uin, msgId, gameSession.GetCurrentS2CMessageId())
	// 	log.Debugf("HandleMessageReceiptSend s2cMsgId not match uin=%d,clientId=%d,serverid=%d\n", gameSession.Uin, msgId, gameSession.GetCurrentS2CMessageId())
	// 	return
	// }
	firstQueueMsg := gameSession.S2CMsgQueue.Front()
	if firstQueueMsg == nil {
		log.Warnf("uin=%d,收到的回执消息id 但 未处理的队列消息为空ip=%s", gameSession.Uin, gameSession.LinkSession.Conn().RemoteAddr().String())
		return
	}

	if gameSession.IsNetStatusAuthing() { // 收到KeepSocketAuthRecv的回执
		for e := gameSession.S2CMsgQueue.Front(); e != nil; e = e.Next() {
			s2cAbsMsgEntity := e.Value.(pub.EntityABSMessage)
			if s2cAbsMsgEntity.MessageId == msgId &&
				pub.GetRealSendActionType(s2cAbsMsgEntity.EntityMessagePairList[0].ActionType) != int32(pf.PFActionType_KeepSocketAuthRecv) {
				gameSession.S2CMsgQueue.Remove(e)
				gameSession.SetNetStatusSendable()
				// hasMatchMsgId = true
				break
			}
		}
		return
	}

	for e := gameSession.S2CMsgQueue.Front(); e != nil; e = e.Next() {
		s2cAbsMsgEntity := e.Value.(pub.EntityABSMessage)
		if s2cAbsMsgEntity.MessageId == msgId {
			gameSession.S2CMsgQueue.Remove(e)
			break
		}
	}
	return
}
