package logic

import (
	"github.com/gogo/protobuf/proto"
	"time"
	"zerogame.info/profile"
	"zerogame.info/taphero/defs"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service"
)

func HandleGetPayToken(gameSession *entity.GameSession, now time.Time) (messages pub.EntityMessagePairList) {
	return
}
func HandleOrderApply(gameSession *entity.GameSession, orderIdList []string, now time.Time) (messages pub.EntityMessagePairList) {
	if len(orderIdList) == 0 {
		return
	}

	user, ok := service.GetUserService().Get(gameSession.Uin, now)
	if !ok {
		return
	}
	if user.GetChannel() == defs.CHANNEL_APPSTORE {
		messages = handleAppStoreOrderApply(gameSession, &user, orderIdList, now)
		return

	}
	if user.GetChannel() == defs.CHANNEL_GOOGLEPLAY {
		messages = handleGooglePlayOrderApply(gameSession, &user, orderIdList, now)
		return
	}

	messages = handleOrderApply(gameSession, &user, orderIdList, now)
	return
}

func getPayStatus(orderIdList []string, orderMap profile.PayOrderMap) (pfList []*pf.OrderStatus) {
	for _, orderId := range orderIdList {
		if _, ok := orderMap[orderId]; ok {
			order := orderMap[orderId]
			if order.IsStatusUnHandled() {
				pf := &pf.OrderStatus{OrderID: proto.String(orderId), Status: proto.Int32(SUCC_STATUS)}
				pfList = append(pfList, pf)
			}
		}
	}
	return
}
