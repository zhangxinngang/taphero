package net

import (
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"time"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
)

func DecodeABSMessage(msgs []byte) (absMsg pub.ABSMessage, ok bool) {
	if msgs == nil {
		ok = false
		return
	}
	absMsg = pub.ABSMessage{}
	err := proto.Unmarshal(msgs, &absMsg)
	if err != nil {
		ok = false
	}
	ok = true

	return
}

func EncodeABSMessageWithHeader(recvABSMsgEntity pub.EntityABSMessage, now time.Time) []byte {
	if recvABSMsgEntity.Time == 0 {
		recvABSMsgEntity.Time = int32(now.Unix())
	}

	msgList := make([]pub.MessagePair, len(recvABSMsgEntity.EntityMessagePairList))
	pubABSmessage := recvABSMsgEntity.ToPB(msgList)
	message, _ := proto.Marshal(&pubABSmessage)

	result := make([]byte, len(message)+4)
	binary.BigEndian.PutUint32(result[0:4], uint32(len(message)))
	copy(result[4:], message)
	return result
}
func GetActionTypeName(actionType int32) string {
	return pf.PFActionType_PFActionTypeDetail_name[actionType]

}
