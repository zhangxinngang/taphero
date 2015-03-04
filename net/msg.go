package net

import (
	"github.com/gogo/protobuf/proto"
	"zerogame.info/taphero/defs"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service/user_add_attr"
)

func EncodeMessageReceiptRecv(msgId int32) (messages pub.EntityMessagePairList) {
	m := pf.MessageReceiptRecv{ReadMessageId: proto.Int32(msgId)}
	message_bytes, _ := proto.Marshal(&m)

	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_MessageReceiptRecv), message_bytes)}
	return

}

func EncodeGameDataBaseRecv(userAddAttr user_add_attr.UserAddAttr) (messages pub.EntityMessagePairList) {
	m := pf.GameDataBaseRecv{
		LastLogoutTM: proto.Int32(int32(userAddAttr.GetLastOffTime().Unix())),
		Energy:       proto.Int32(userAddAttr.GetEnergy()),
		EnergyAddTM:  proto.Int32(int32(userAddAttr.GetEnergyTime().Unix())),
	}
	if userAddAttr.GetLastOffTime() == defs.UninitedTime {
		m.LastLogoutTM = proto.Int32(0)
	}

	message_bytes, _ := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_GameDataBase), message_bytes),
	}
	return
}
func EncodeSyncLogoutTimeRecv() (messages pub.EntityMessagePairList) {
	// m := pf.SyncLogoutTimeRecv{}
	// message_bytes, err := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_SyncLogoutTime), nil),
	}
	return
}

func EncodeKeepSocketAuthRecv(nextMessageId int32) (messages pub.EntityMessagePairList) {
	m := pf.KeepSocketAuthRecv{NextMessageId: proto.Int32(nextMessageId)}
	message_bytes, _ := proto.Marshal(&m)

	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_KeepSocketAuthRecv), message_bytes)}
	return

}
func EncodeCheckVersionRecv(isVersionEnabled int32, adsEnabled int32) (messages pub.EntityMessagePairList) {
	m := pf.CheckVersionRecv{
		VerEnabled: proto.Int32(isVersionEnabled),
		AdsEnabled: proto.Int32(adsEnabled),
	}
	message_bytes, _ := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_CheckVersion), message_bytes),
	}
	return
}
func EncodeOrderApplyRecv(statuses []*pf.OrderStatus) (messages pub.EntityMessagePairList) {
	m := pf.OrderApplyRecv{OrderList: statuses}
	message_bytes, _ := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_OrderApply), message_bytes),
	}
	return
}
func EncodeGotoDungeonRecv(status int32, userAddAttr *user_add_attr.UserAddAttr) (messages pub.EntityMessagePairList) {
	m := pf.GotoDungeonRecv{
		Status:      proto.Int32(status),
		Energy:      proto.Int32(userAddAttr.GetEnergy()),
		EnergyAddTM: proto.Int32(int32(userAddAttr.GetEnergyTime().Unix())),
	}
	message_bytes, _ := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_GotoDungeon), message_bytes),
	}
	return
}
func EncodeBuyEnergyRecv(userAddAttr *user_add_attr.UserAddAttr) (messages pub.EntityMessagePairList) {
	m := pf.BuyEnergyRecv{
		Energy:      proto.Int32(userAddAttr.GetEnergy()),
		EnergyAddTM: proto.Int32(int32(userAddAttr.GetEnergyTime().Unix())),
	}
	message_bytes, _ := proto.Marshal(&m)
	messages = pub.EntityMessagePairList{
		pub.PackRecvMsgBody(int32(pf.PFActionType_BuyEnergy), message_bytes),
	}
	return
}
