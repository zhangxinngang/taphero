package logic

import (
	"encoding/json"
	"fmt"
	"github.com/0studio/goauth/utils"
	"github.com/gogo/protobuf/proto"
	//"strings"
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

/*
	Google Play 支付查询
	https://www.googleapis.com/androidpublisher/v2/applications/%s/purchases/products/%s/tokens/%s
*/
const (
	PAY_GOOGLE_PLAY_TIMEOUT = 30 * 1000 // 10s
)

func getGooglePlayResponseTry5Times(url string, now time.Time) (response []byte, err error) {
	for i := 0; i < 5; i++ {
		response, err = utils.GetHttpResponseAsJson(url, now, PAY_GOOGLE_PLAY_TIMEOUT)
		if err == nil {
			return
		}
	}
	return
}

// """
/*
{
 "error": {
  "errors": [
   {
    "domain": "global",
    "reason": "required",
    "message": "Login Required",
    "locationType": "header",
    "location": "Authorization"
   }
  ],
  "code": 401,
  "message": "Login Required"
 }
}
*/

type GooglePlayResponse struct {
	Error            GooglePlayErrorResponse
	PurchaseState    int32  `json:"purchaseState,omitempty"`
	DeveloperPayLoad string `json:"developerPayload,omitempty"`
}

type GooglePlayErrorResponse struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type GooglePlayParams struct {
	AppName   string `json:"applications,omitempty"`
	ProductId string `json:"products,omitempty"`
	token     string `json:"token,omitempty"`
}

func GetGooglePlayResponse(content string, now time.Time) (jsonData GooglePlayResponse, err error) {
	gpp := GooglePlayParams{}
	json.Unmarshal([]byte(content), &gpp)
	/*
		strs := strings.SplitN(content, "-", -1)
		appname := strs[0]
		productid := strs[1]
		token := strs[2]
	*/
	fmt.Println(content, "content")
	url := fmt.Sprintf("https://www.googleapis.com/androidpublisher/v2/applications/%s/purchases/products/%s/tokens/%s", gpp.AppName, gpp.ProductId, gpp.token)
	response, err := getGooglePlayResponseTry5Times(url, now)
	fmt.Println(string(response), "response...")
	if err != nil {
		return
	}
	jsonData = GooglePlayResponse{}
	json.Unmarshal(response, &jsonData)

	return
}

func handleGooglePlayOrderApply(gameSession *entity.GameSession, user *profile.User, orderIdList []string, now time.Time) (messages pub.EntityMessagePairList) {
	var pfStatusList []*pf.OrderStatus
	for _, clientInfo := range orderIdList {
		order, ok := handleOrderApplyGooglePlay(user, clientInfo, now)
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

func handleOrderApplyGooglePlay(user *profile.User, clientInfo string, now time.Time) (payOrder profile.PayOrder, succ bool) {
	gpp := GooglePlayParams{}
	json.Unmarshal([]byte(clientInfo), &gpp)
	Uin := user.GetUin()

	jsonData, err := GetGooglePlayResponse(clientInfo, now)
	if err != nil {
		log.Error("google play get pay response error:", err)
		return
	}
	if jsonData.Error.Code != 0 {
		log.Errorf("google play err: message:%s ,code:%d", jsonData.Error.Message, jsonData.Error.Code)
		return
	}
	if jsonData.PurchaseState != 0 || jsonData.DeveloperPayLoad != fmt.Sprintf("DeveloperPayloadITEM%s", gpp.ProductId) || jsonData.DeveloperPayLoad == "" {
		log.Errorf("google play err: purchase info is wrong PurchaseState = %d and DeveloperPayLoad = %s", jsonData.PurchaseState, jsonData.DeveloperPayLoad)
		return
	}

	var (
		ok bool
	)

	payOrder, ok = service.GetPayOrderService().Get(Uin, gpp.ProductId)
	if ok && payOrder.IsStatusHandled() {
		log.Infof("already handled pay payOrder,Uin=%d,orderid=%s", Uin, gpp.ProductId)
		return
	}

	payOrder = profile.PayOrder{}
	payOrder.SetUin(uint64(Uin))
	payOrder.SetAccountId(user.GetAccountId())
	payOrder.SetOrderId(clientInfo)
	payOrder.SetProductId(gpp.ProductId)
	payOrder.SetChannel(user.GetChannel())
	payOrder.SetServerId(conf.GetServer())
	payOrder.SetStatusHandled()
	payOrder.SetRecvData(clientInfo)
	payOrder.SetMoney(0)
	payOrder.SetOrderType(0)
	payOrder.SetCreateTime(now)
	succ = service.GetPayOrderService().Add(payOrder)

	if !succ {
		log.Errorf("google play insert payorder fail Uin=%d,OrderId=%s,ProductId=%d,recv_data=%s",
			Uin, payOrder.GetOrderId(), payOrder.GetProductId(), payOrder.GetRecvData())
		return

	}
	log.Infof("google play appstore pay payOrder succ ,Uin=%d,orderid=%s", user.GetUin(), payOrder.GetOrderId())

	return
}
