package pub

import (
	"github.com/gogo/protobuf/proto"
	"time"
)

type EntityMessagePair struct {
	ActionType   int32
	MessageBytes []byte
}

func (Msg *EntityMessagePair) ToPB(pb *MessagePair) {
	pb.ActionType = proto.Int32(Msg.ActionType)
	pb.MessageBytes = Msg.MessageBytes
}

type EntityMessagePairList []EntityMessagePair

func (MsgList *EntityMessagePairList) ToPB(pb []MessagePair) {
	for idx, _ := range *MsgList {
		(*MsgList)[idx].ToPB(&(pb[idx]))
	}
	return
}

type EntityABSMessage struct {
	EntityMessagePairList EntityMessagePairList
	Token                 string
	Subversion            int32
	MessageId             int32
	Time                  int32 // unix timestamp
}

func (absMsg *EntityABSMessage) ToPB(msgList []MessagePair) (msg ABSMessage) {
	absMsg.EntityMessagePairList.ToPB(msgList)
	msgPointerList := make([]*MessagePair, len(absMsg.EntityMessagePairList))
	for idx, _ := range msgList {
		msgPointerList[idx] = &(msgList[idx])
	}
	msg = ABSMessage{
		MsgList:    msgPointerList,
		Timestamp:  proto.Int32(absMsg.Time),
		Subversion: proto.Int32(absMsg.Subversion),
		MessageId:  proto.Int32(absMsg.MessageId),
		Token:      proto.String(absMsg.Token),
	}
	return
}

func PackRecvMsgBody(action_type int32, message_bytes []byte) (msg EntityMessagePair) {
	action_type += RECV_SHIFT
	msg = EntityMessagePair{ActionType: action_type, MessageBytes: message_bytes}
	return
}

func NewMsg(msgId int32, msgs EntityMessagePairList, token string, subVersion int32, now time.Time) EntityABSMessage {
	// Subversion: version.GetCurrentServerVersion
	return EntityABSMessage{
		Time:                  int32(now.Unix()),
		MessageId:             msgId,
		EntityMessagePairList: msgs,
		Token:      token,
		Subversion: subVersion,
	}
}

var (
	EmptyABSMessage = EntityMessagePair{
		ActionType:   0,
		MessageBytes: []byte{}}
	EmptyABSMessageArr = []EntityMessagePair{EmptyABSMessage}
)
