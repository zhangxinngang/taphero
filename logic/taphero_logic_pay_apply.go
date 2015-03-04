package logic

import (
	"github.com/gogo/protobuf/proto"
	"time"
	"zerogame.info/profile"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service"
)

func handleOrderApply(gameSession *entity.GameSession, user *profile.User, orderIdList []string, now time.Time) (messages pub.EntityMessagePairList) {
	var pfStatusList []*pf.OrderStatus
	for _, orderId := range orderIdList {
		order, ok := updateOrderStatusHandledById(user, orderId, now)
		if ok {
			pf := &pf.OrderStatus{
				OrderID:    proto.String(order.GetOrderId()),
				Status:     proto.Int32(SUCC_STATUS),
				AppStoreID: proto.String(order.GetProductId()),
			}
			pfStatusList = append(pfStatusList, pf)
		}
	}

	if len(pfStatusList) == 0 {
		return
	}
	messages = net.EncodeOrderApplyRecv(pfStatusList)
	return
}

func handleUnhandledPayOrder(user *profile.User, now time.Time) (messages pub.EntityMessagePairList) {
	log.Debug("handleUnhandledPayOrder", user.GetUin())
	//处理漏单
	payOrderList := service.GetPayOrderService().GetAllUnhandledOrder(user.GetUin())
	var pfStatusList []*pf.OrderStatus
	if len(payOrderList) != 0 {
		for idx, _ := range payOrderList {
			ok := updateOrderStatusHandled(user, &(payOrderList[idx]), now)
			if ok {
				pf := &pf.OrderStatus{
					OrderID:    proto.String(payOrderList[idx].GetOrderId()),
					Status:     proto.Int32(SUCC_STATUS),
					AppStoreID: proto.String(payOrderList[idx].GetProductId()),
				}
				pfStatusList = append(pfStatusList, pf)
			}
		}
	} else {
		return
	}
	if len(pfStatusList) == 0 {
		return
	}
	messages = net.EncodeOrderApplyRecv(pfStatusList)
	return
}

func updateOrderStatusHandled(user *profile.User, payOrder *profile.PayOrder, now time.Time) (succ bool) {
	payOrder.SetStatusHandled()
	succ = service.GetPayOrderService().Set(*payOrder)
	if !succ {
		log.Errorf("insert payorder fail Uin=%d,OrderId=%s,ProductId=%d,recv_data=%s",
			user.GetUin(), payOrder.GetOrderId(), payOrder.GetProductId(), payOrder.GetRecvData())
		return
	}
	log.Infof("apply pay payOrder succ ,Uin=%d,orderid=%s", user.GetUin(), payOrder.GetOrderId())

	return
}
func updateOrderStatusHandledById(user *profile.User, orderId string, now time.Time) (payOrder profile.PayOrder, succ bool) {
	payOrder, ok := service.GetPayOrderService().Get(user.GetUin(), orderId)

	if ok && payOrder.IsStatusHandled() {
		log.Infof("already handled pay payOrder,Uin=%d,orderid=%s", user.GetUin(), orderId)
		return
	}
	if !ok {
		return
	}
	succ = updateOrderStatusHandled(user, &payOrder, now)
	return
}
