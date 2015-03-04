package pub

import (
	"github.com/gogo/protobuf/proto"
)

const SEND_SHIFT = 10000
const RECV_SHIFT = 20000

// func PackSendABS(action_type int32, message_bytes []byte, token string) (outbuf []byte) {
// 	action_type += SEND_SHIFT
// 	msg := &ABSMessage{
// 		MsgList: []*MessagePair{&MessagePair{ActionType: proto.Int32(action_type), MessageBytes: message_bytes}},
// 		// Subversion: proto.Int32(version.GetCurrentServerVersion()),
// 		Token: proto.String(token)}
// 	outbuf, _ = proto.Marshal(msg)
// 	return
// }
// func (abs *ABSMessage) SetMessageId(msgId int32) {
// 	abs.MessageId = proto.Int32(msgId)
// }

func GetRealSendActionType(streamActionType int32) int32 {
	return streamActionType - SEND_SHIFT
}
func GetRealRecvActionType(protoActionType int32) int32 {
	return protoActionType - RECV_SHIFT
}

func GetClientRecvActionType(protoActionType int32) int32 {
	return protoActionType + RECV_SHIFT
}

func GetClientSendActionType(protoActionType int32) int32 {
	return protoActionType + SEND_SHIFT
}

func DecodeABSMessage(msgs []byte) (absMsg ABSMessage, ok bool) {
	if msgs == nil {
		ok = false
		return
	}
	absMsg = ABSMessage{}
	err := proto.Unmarshal(msgs, &absMsg)
	if err != nil {
		ok = false
	}
	ok = true

	return
}
