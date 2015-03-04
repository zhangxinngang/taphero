package logic

import (
	"fmt"
	key "github.com/0studio/storage_key"
	"github.com/fanngyuan/link"
	"time"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/resource"
	"zerogame.info/taphero/service"
	"zerogame.info/taphero/timer"
)

const (
	MSG_CHECKER_DUR_SECONDS = 30
)

func StartPlayerSessionScheduler(uin key.KeyUint64) {
	playerSessionCheckerGoRoutine := func() {
		for {
			gameSession, ok := resource.GetGameSession(uin)
			if !ok {
				return
			}

			select {
			case <-timer.GetSessionTimeWheelChan(time.Second * MSG_CHECKER_DUR_SECONDS):
				if !doSendMsgChecker(uin) {
					CloseSession(uin)
					break
				}

			case msgHandle := <-gameSession.MsgHandleChan:
				msgHandle.GameSession = &gameSession
				HandleUserMsgList(&msgHandle)
				resource.PutGameSession(gameSession.Uin, gameSession)
				// msgHandle.RecvMsgHandleChan <- msgHandle
			case ackMsgHandle := <-gameSession.AckMsgHandleChan:
				if SendAck(ackMsgHandle.C2SABSMsg, gameSession.LinkSession, ackMsgHandle.Now) {
					ackMsgHandle.GameSession = &gameSession
					HandleUserMsgList(&ackMsgHandle)
					resource.PutGameSession(gameSession.Uin, *ackMsgHandle.GameSession)
				}

			case s2cAbsMsgEntity := <-gameSession.S2CEntityABSMessagePusher:
				gameSession.S2CMsgQueue.PushBack(s2cAbsMsgEntity)
				resource.PutGameSession(gameSession.Uin, gameSession)
			// case inscrC2SMessageId := <-gameSession.C2SMessageIdInscr:
			// 	if inscrC2SMessageId == 1 {
			// 		// wonot go herer
			// 		gameSession.SetNextC2SMessageId(1)
			// 	} else {
			// 		log.Debug("old gameSession.GetNextC2SMessageId() ", gameSession.GetNextC2SMessageId())
			// 		gameSession.SetNextC2SMessageId(gameSession.GetNextC2SMessageId() + 1)
			// 		log.Debugf("insc c2s_msgid uin=%d new c2smsgid=%d", gameSession.Uin, gameSession.GetNextC2SMessageId())
			// 	}
			// 	resource.PutGameSession(gameSession.Uin, gameSession)
			case gameSessionReconnStruct := <-gameSession.GameSessionReconnStruct:
				log.Debugf("session_scheduler_处理再次auth,uin=%d,newC2sNextMessageId=%d", gameSession.Uin, gameSessionReconnStruct.C2SNextMessageId)
				gameSession.LinkSession = gameSessionReconnStruct.Session
				gameSession.SetNetStatusAuthing()
				gameSession.SetNextC2SMessageId(gameSessionReconnStruct.C2SNextMessageId) // 如果客户端给我发auth协议则重置为1
				if gameSessionReconnStruct.C2SNextMessageId == 1 {                        // 如果client 要求重置id 为1 ,则清空我的消息队列
					for e := gameSession.S2CMsgQueue.Front(); e != nil; e = e.Next() {
						gameSession.S2CMsgQueue.Remove(e)
					}
				}
				resource.PutGameSession(gameSession.Uin, gameSession)
				resource.PutUin(gameSessionReconnStruct.Session.Id(), gameSession.Uin)
				// 重连上来 ，马上把队列中的消息刷到客户端
				doSendMsgChecker(gameSession.Uin)
			case msg := <-gameSession.MsgSender:
				gameSession.SendMsgSelf(msg.Msg, msg.Now)
				resource.PutGameSession(gameSession.Uin, gameSession)
			case <-gameSession.CheckOrder:
				user, ok := service.GetUserService().Get(gameSession.Uin, time.Now())
				if !ok {
					return
				}
				msg := handleUnhandledPayOrder(&user, time.Now())
				gameSession.SendMsg(msg, time.Now())
			case closedTime := <-gameSession.OnTcpClosed:
				doOnTcpClosed(&gameSession, closedTime, time.Now())
			}
		}
	}
	go playerSessionCheckerGoRoutine()
}
func doSendMsgChecker(uin key.KeyUint64) bool {
	log.Debug("doSendMsgChecker", uin)
	gameSession, ok := resource.GetGameSession(uin)
	if !ok {
		return false
	}
	now := time.Now()
	if gameSession.LinkSession.IsClosed() {
		// 当前tcp 已断， 等待新的tcp 连上
		return true
	}
	if gameSession.IsNetStatusDiabled() {
		return true

	}
	for e := gameSession.S2CMsgQueue.Front(); e != nil; e = e.Next() {
		s2cAbsMsgEntity := e.Value.(pub.EntityABSMessage)
		// for s2cMessageId, s2cAbsMsgEntity := range gameSession.S2CMsgQueue. {
		if now.Sub(time.Unix(int64(s2cAbsMsgEntity.Time), 0)).Seconds() < MSG_CHECKER_DUR_SECONDS {
			log.Debug("wonot resend_msg_for_dur_too_short", uin, s2cAbsMsgEntity.MessageId)
			continue
		}
		if gameSession.IsNetStatusAuthing() {
			// 如果是在等待客户端发回给我KeepSocketAuthRecv的回执的过程中， 那么我只能向客户端发送KeepSocketAuthRecv，以保证其一定能收到KeepSocketAuthRecv，以保证其一定能收到
			// 且在收到回执之前不向客户端重传其他任何消息
			if pub.GetRealSendActionType(s2cAbsMsgEntity.EntityMessagePairList[0].ActionType) != int32(pf.PFActionType_KeepSocketAuthRecv) {
				continue
			}
			log.Debugf("resend_msg,uin=%d,sessionid=%d,s2cMsgId=%d,msg=%v", uin, gameSession.LinkSession.Id(), s2cAbsMsgEntity.MessageId, s2cAbsMsgEntity)
			gameSession.LinkSession.SendPacket(net.WritePacket(net.EncodeABSMessageWithHeader(s2cAbsMsgEntity, now)))
			return true
		}

		log.Debugf("resend_msg uin=%d,sessionid=%d,s2cMsgId=%d,msg=%v", uin, gameSession.LinkSession.Id(), s2cAbsMsgEntity.MessageId, s2cAbsMsgEntity)
		gameSession.LinkSession.SendPacket(net.WritePacket(net.EncodeABSMessageWithHeader(s2cAbsMsgEntity, now)))
	}
	return true
}

func CloseSession(uin key.KeyUint64) {
	resource.DelGameSession(uin)
}
func OnTcpClosed(session link.SessionAble, closedTime time.Time, reason interface{}) {
	fmt.Println(reason)
	gameSession, ok := resource.GetGameSessionByTcpSessionId(session.Id())
	if !ok {
		return
	}
	gameSession.TryOnTcpClosed(closedTime)
}
func doOnTcpClosed(gameSession *entity.GameSession, closedTime time.Time, now time.Time) {
	onTcpClosed4LastOffLineTime(gameSession, closedTime, now)
}
