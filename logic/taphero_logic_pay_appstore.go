package logic

import (
	"encoding/json"
	"github.com/0studio/goauth/utils"
	"github.com/gogo/protobuf/proto"
	"strings"
	"time"
	"zerogame.info/profile"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/entity"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/net"
	"zerogame.info/taphero/pf"
	"zerogame.info/taphero/pub"
	"zerogame.info/taphero/service"
)

//
//
// func url() string {
// 	if conf.IsModePro() {
// 		return "https://buy.itunes.apple.com/verifyReceipt"
// 	}
// 	return "https://sandbox.itunes.apple.com/verifyReceipt"

// }
const (
	PAY_APP_STORE_TIMEOUT = 30 * 1000 // 10s

)

func getAppStoreResponseTry5Times(url string, now time.Time, content []byte) (response []byte, err error) {
	for i := 0; i < 5; i++ {
		response, err = utils.PostHttpResponse(url, content, now, PAY_APP_STORE_TIMEOUT)
		if err == nil {
			return
		}
	}
	return
}

type AppStorePayResponse struct {
	Status    int32
	OrderType int32            // 0普通定单，1沙盒测试定单
	Info      AppStoreRecvData `json:"receipt,omitempty"`
}
type AppStoreRecvData struct {
	TransactionId string `json:"transaction_id,omitempty"`
	ProductId     string `json:"product_id,omitempty"`
}

// """
// {"status":21002, "exception":"java.lang.ClassCastException"}
//
// {"receipt":{"original_purchase_date_pst":"2012-11-12 03:10:48 America/Los_Angeles",
// "purchase_date_ms":"1352718648392",
// "unique_identifier":"98de2b8cb8b973773538c5c8743e1043677b9201",
// "original_transaction_id":"1000000058437063",
// "bvrs":"50000",
// "transaction_id":"1000000058437063",
// "quantity":"1",
// "unique_vendor_identifier":"FB22FBAA-CFC9-4ABD-9522-BCECF78C2866",
// "item_id":"577373064",
// "product_id":"com.yyshtech.gold_6",
// "purchase_date":"2012-11-12 11:10:48 Etc/GMT",
// "original_purchase_date":"2012-11-12 11:10:48 Etc/GMT",
// "purchase_date_pst":"2012-11-12 03:10:48 America/Los_Angeles",
// "bid":"com.yyshtech.zhajinhua", "original_purchase_date_ms":"1352718648392"}, "status":0}
// product_id: com.yysh.goldcount.6

// How do I verify my receipt (iOS)?
// Always verify your receipt first with the production URL; proceed to verify with the sandbox URL if you receive a 21007 status code. Following this approach ensures that you do not have to switch

const (
	// 0普通定单，1沙盒测试定单
	PAY_ORDER_TYPE_COMMON  = 0
	PAY_ORDER_TYPE_SANDBOX = 1
)

func GetAppStoreResponse(content []byte, now time.Time) (jsonData AppStorePayResponse, err error) {
	response, err := getAppStoreResponseTry5Times("https://buy.itunes.apple.com/verifyReceipt", now, content)
	// response, err := getAppStoreResponseTry5Times("https://sandbox.itunes.apple.com/verifyReceipt", content)

	if err != nil {
		return
	}
	jsonData = AppStorePayResponse{}
	json.Unmarshal(response, &jsonData)
	jsonData.OrderType = PAY_ORDER_TYPE_COMMON
	if jsonData.Status == 21007 {
		jsonData.OrderType = PAY_ORDER_TYPE_SANDBOX
		response, err = getAppStoreResponseTry5Times("https://sandbox.itunes.apple.com/verifyReceipt", now, content)
		if err != nil {
			return
		}
		json.Unmarshal(response, &jsonData)

	}

	return
}

type ClientAppStoreSend struct {
	TransactionId string `json:"transaction_id,omitempty"`
	ReceiptData   string `json:"receipt-data,omitempty"`
}

func handleAppStoreOrderApply(gameSession *entity.GameSession, user *profile.User, orderIdList []string, now time.Time) (messages pub.EntityMessagePairList) {
	var pfStatusList []*pf.OrderStatus
	for _, clientInfo := range orderIdList {
		order, ok := handleOrderApplyAppstore(user, clientInfo, now)
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

// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html#//apple_ref/doc/uid/TP40010573-CH104-SW1
func handleOrderApplyAppstore(user *profile.User, clientInfo string, now time.Time) (payOrder profile.PayOrder, succ bool) {
	send := ClientAppStoreSend{}
	Uin := user.GetUin()
	err := json.Unmarshal([]byte(clientInfo), &send)
	if err != nil {
		// ban.AddDefaultBinUin(user.Uin, "appstore pay error")
		log.Errorf("appstore pay wrong format orderId Uin=%d,send.TransactionId=%s", Uin, clientInfo)
		return
	}
	log.Infof("appstore pay Uin=%d,send.TransactionId=%s", Uin, send.TransactionId)
	if send.ReceiptData == "" {
		return
	}

	if len(send.TransactionId) >= 19 || len(send.TransactionId) <= 8 {
		log.Errorf("appstore pay wrong format transactionid Uin=%d,send.TransactionId=%s", Uin, send.TransactionId)
		return
	}
	if strings.Contains(send.TransactionId, "-") {
		log.Errorf("appstore pay wrong format transactionid Uin=%d,send.TransactionId=%s", Uin, send.TransactionId)
		return

	}
	if strings.Contains(send.TransactionId, "com.urus.iap") {
		log.Errorf("appstore pay wrong format transactionid Uin=%d,send.TransactionId=%s", Uin, send.TransactionId)
		return
	}

	jsonData, err := GetAppStoreResponse([]byte(clientInfo), now)
	if err != nil {
		log.Error("appstore get pay response error:", err)
		return
	}
	if jsonData.Status != 0 {
		log.Error("appstore get pay response status !=0:", user.GetUin(), err, jsonData.Status, jsonData, clientInfo)
		return
	}

	var (
		ok bool
	)

	// dao.PayOrder(Uin, send.TransactionId)
	payOrder, ok = service.GetPayOrderService().Get(Uin, send.TransactionId)
	if ok && payOrder.IsStatusHandled() {
		log.Infof("already handled pay payOrder,Uin=%d,orderid=%s", Uin, send.TransactionId)
		return
	}

	// if ok && payOrder.IsStatusUnHandled() {
	// 	payOrder.SetStatusHandled()
	// 	succ = service.GetPayOrderService().UpdateStatus(payOrder)
	// } else {
	payOrder = profile.PayOrder{}
	payOrder.SetUin(uint64(Uin))
	payOrder.SetAccountId(user.GetAccountId())
	payOrder.SetOrderId(send.TransactionId)
	payOrder.SetProductId(jsonData.Info.ProductId)
	payOrder.SetChannel(user.GetChannel())
	payOrder.SetServerId(conf.GetServer())
	payOrder.SetStatusHandled()
	payOrder.SetRecvData(clientInfo)
	payOrder.SetMoney(0)
	payOrder.SetOrderType(jsonData.OrderType)
	payOrder.SetCreateTime(now)
	succ = service.GetPayOrderService().Add(payOrder)
	// }

	if !succ {
		log.Errorf("appstore insert payorder fail Uin=%d,OrderId=%s,ProductId=%d,recv_data=%s",
			Uin, payOrder.GetOrderId(), payOrder.GetProductId(), payOrder.GetRecvData())
		return

	}
	log.Infof("apply appstore pay payOrder succ ,Uin=%d,orderid=%s", user.GetUin(), payOrder.GetOrderId())

	return
}
