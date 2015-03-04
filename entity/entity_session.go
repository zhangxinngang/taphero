package entity

import (
	"github.com/0studio/storage_key"
	"github.com/fanngyuan/link"
	"time"
	// "zerogame.info/profile/auth"
	"container/list"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
)

// session {2256197860196353 0xc208378000 map[] 0 0 0xc208341800 0 <nil> <nil> <nil>}
type MsgHandle struct {
	SDKInfo interface{}

	GameSession *GameSession
	C2SABSMsg   pub.ABSMessage
	// S2CABSMsgEntity   pub.EntityABSMessage
	Now time.Time
	// HandleFunc func(*MsgHandle)
	// RecvMsgHandleChan chan MsgHandle
}

type MsgPack struct {
	Msg pub.EntityMessagePairList
	Now time.Time
}

// type MsgMap map[int32]pub.EntityABSMessage // key msgId
type GameSession struct {
	Uin                       key.KeyUint64
	LinkSession               link.SessionAble
	S2CMsgQueue               *list.List // 我需要保证一定要发给客户端的消息 需要客户端确认收到的消息列表
	currentS2CMessageId       int32      //自增序列，每发一条消息,++
	nextC2SMessageId          int32      // client to server msgId
	MsgHandleChan             chan MsgHandle
	AckMsgHandleChan          chan MsgHandle
	netStatus                 int8 // 0,not_sendable(不可向客户端发包),1,authing(正在认证),2ok（可以客户端发包）
	S2CEntityABSMessagePusher chan pub.EntityABSMessage
	MsgSender                 chan MsgPack
	C2SMessageIdInscr         chan int
	GameSessionReconnStruct   chan GameSessionReconnStruct
	CheckOrder                chan bool
	OnTcpClosed               chan time.Time
}

func NewGameSession(Uin key.KeyUint64, session link.SessionAble) (gameSession GameSession) {
	gameSession.Uin = Uin
	gameSession.LinkSession = session
	gameSession.S2CMsgQueue = list.New()
	gameSession.MsgHandleChan = make(chan MsgHandle)
	gameSession.AckMsgHandleChan = make(chan MsgHandle)
	gameSession.GameSessionReconnStruct = make(chan GameSessionReconnStruct)
	gameSession.C2SMessageIdInscr = make(chan int)
	gameSession.S2CEntityABSMessagePusher = make(chan pub.EntityABSMessage)
	gameSession.MsgSender = make(chan MsgPack)
	gameSession.currentS2CMessageId = 1
	gameSession.CheckOrder = make(chan bool)
	gameSession.OnTcpClosed = make(chan time.Time)
	// gameSession.nextC2SMessageId = 1
	gameSession.SetNetStatusAuthing()
	return

}

func (session *GameSession) TryCheckPayOrder() {
	session.CheckOrder <- true

}
func (session *GameSession) TryOnTcpClosed(closedTime time.Time) {
	session.OnTcpClosed <- closedTime
}

func (session GameSession) IsNetStatusSendable() bool { // 认证通过，且客户端明确收到KeepSocketAuthRecv（发送给我过回执）
	// 服务器可不可以给客户端发普通的消息了，
	// 只有当客户端给我发了KeepSocketAuthSend,我回复了此消息的回 执
	// 并且我全客户端发送了KeepSocketAuthRecv,且客户端回复我我回执
	// 此时才真正认为 相互可以发消息了，
	// 如果当发现tcp 连接突然断开后， 需要将此变成false ,
	// 需要重新走上面的步骤
	return session.netStatus == 2
}

func (session GameSession) IsNetStatusAuthing() bool { // 正在认证阶段（已回复KeepSocketAuthRecv,但是还没收到客户端的回执消息）
	// 服务器可不可以给客户端发普通的消息了，
	// 只有当客户端给我发了KeepSocketAuthSend,我回复了此消息的回 执
	// 并且我全客户端发送了KeepSocketAuthRecv,且客户端回复我我回执
	// 此时才真正认为 相互可以发消息了，
	// 如果当发现tcp 连接突然断开后， 需要将此变成false ,
	// 需要重新走上面的步骤
	return session.netStatus == 1
}
func (session GameSession) IsNetStatusDiabled() bool {
	return session.netStatus == 0
}

func (session *GameSession) SetNetStatusSendable() {
	session.netStatus = 2
}

func (session *GameSession) SetNetStatusAuthing() {
	session.netStatus = 1
}
func (session *GameSession) SetNetStatusDisable() {
	session.netStatus = 0
}
func (session *GameSession) PushMsgHandle(msgHandle MsgHandle) {
	// 当时msgId比我期望的msgid大时， server 不立即处理此消息， 而是把它push到gamesessoin的goroutine中， 等待再次判断可不可处理
	go func() {
		session.MsgHandleChan <- msgHandle
	}()
}

func (session *GameSession) PushAckMsgHandle(msgHandle MsgHandle) {
	// 当时msgId比我期望的msgid大时， server 不立即处理此消息， 而是把它push到gamesessoin的goroutine中， 等待再次判断可不可处理
	go func() {
		session.AckMsgHandleChan <- msgHandle
	}()
}

func (session *GameSession) PushS2CMsgToQueue(s2cABSMsgEntity pub.EntityABSMessage) {
	go func() {
		session.S2CEntityABSMessagePusher <- s2cABSMsgEntity
	}()
}

func (session *GameSession) PushS2CMsgToQueueSelf(s2cABSMsgEntity pub.EntityABSMessage) {
	session.S2CMsgQueue.PushBack(s2cABSMsgEntity)
}

func (session GameSession) GetNextC2SMessageId() int32 {
	return session.nextC2SMessageId
}
func (session *GameSession) SetNextC2SMessageId(nextC2SMessageId int32) {
	session.nextC2SMessageId = nextC2SMessageId
}

// func (session *GameSession) InscNextC2SMessageId() {
// 	go func() {
// 		session.C2SMessageIdInscr <- 0 // 发0表示另其自增
// 	}()

// 	// session.nextC2SMessageId = session.nextC2SMessageId + 1
// }

type GameSessionReconnStruct struct {
	C2SNextMessageId int32
	Session          link.SessionAble
}

func (session *GameSession) Reconn(tcpSession link.SessionAble, c2sNextMessageId int32) {
	go func() {
		session.GameSessionReconnStruct <- GameSessionReconnStruct{C2SNextMessageId: c2sNextMessageId, Session: tcpSession}
	}()
}
func (session *GameSession) ResetNextC2SMessageId() {
	// useless
	go func() {
		session.C2SMessageIdInscr <- 1 // 发1表示另其重置为1
	}()

}
func (session GameSession) GetCurrentS2CMessageId() int32 {
	return session.currentS2CMessageId
}
func (session *GameSession) SetCurrentS2CMessageId(nextMessageId int32) {
	session.currentS2CMessageId = nextMessageId
}
func (session *GameSession) SendMsg(msg pub.EntityMessagePairList, now time.Time) { // 只使用其中的msg.EntityMessagePairList 和msg.Time
	go func() {
		session.MsgSender <- MsgPack{msg, now}
	}()
}
func (session *GameSession) SendMsgSelf(msgPairList pub.EntityMessagePairList, now time.Time) {
	if len(msgPairList) == 0 {
		return
	}

	absMsg := pub.NewMsg(session.GetCurrentS2CMessageId(), msgPairList, "", 0, now)
	absMsg.MessageId = session.GetCurrentS2CMessageId() // FIXME: ??? really need

	if pub.GetRealRecvActionType(absMsg.EntityMessagePairList[0].ActionType) == int32(pf.PFActionType_KeepSocketAuthRecv) {
		absMsg.MessageId = 0
	} else {
		session.SetCurrentS2CMessageId(1 + session.GetCurrentS2CMessageId())
		log.Debugf("i_set_s2cmsgid =%d %s", session.GetCurrentS2CMessageId(), session.LinkSession.Conn().RemoteAddr().String())

	}
	data := net.EncodeABSMessageWithHeader(absMsg, time.Now()) //

	session.LinkSession.SendPacket(net.WritePacket(data))
	session.PushS2CMsgToQueueSelf(absMsg)
	log.Debugf("<------s2c[%s]msgid=%d,uin=%d,ip=%s",
		net.GetActionTypeName(pub.GetRealRecvActionType(absMsg.EntityMessagePairList[0].ActionType)),
		absMsg.MessageId,
		session.Uin,
		session.LinkSession.Conn().RemoteAddr().String())

}

// func (session *GameSession) GetSDKAccessCode(channel uint64, now time.Time) string {
// 	// do not call this directly , call resource.GetSDKAccessToken()
// 	if channel == auth.CHANNEL_360 {
// 		sdk360 := session.SDKInfo.(auth.AccessTokenRec360)
// 		if sdk360.ExpireTime.Before(now.Add(time.Second * 60)) {
// 			for i := 0; i < 5; i++ {
// 				if sdk360.Refresh360Token(now) {
// 					session.SDKInfo = sdk360
// 					break
// 				}
// 			}
// 		}
// 		return sdk360.Access_token
// 	}
// 	return ""
// }

type GameSessionMap map[key.KeyUint64]GameSession
