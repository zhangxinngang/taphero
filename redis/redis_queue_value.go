package redis

import (
	"time"
)

type RedisQueueValue struct {
	Type                  int32
	Platform              uint64   `json:"platform,omitempty"`
	Uin                   string   `json:"uin,omitempty"`
	SUin                  string   `json:"suin,omitempty"`
	ProcessList           []uint64 `json:"proc,omitempty"`
	Time                  int64    `json:"time,omitempty"`
	MyAccountId           string   `json:"myaccountid,omitempty"`
	MyUuid                string   `json:"uuid,omitempty"`
	AccountId             string   `json:"accountid,omitempty"`
	AccountName           string   `json:"accountname,omitempty"`
	Channel               uint64   `json:"chan,omitempty"`
	SelfProcessId         uint64   `json:"pid,omitempty"`
	IsSelfProcessIdAccept int8     `json:"is_pid_accept,omitempty"`
}

func (value RedisQueueValue) IsProcessIdAccept(processId uint64) bool {
	// 判断 此条消息发送者需不需要自己处理自己发送的消息
	if value.IsSelfProcessIdAccept == 1 && value.SelfProcessId == processId {
		return true
	}
	return false

}
func (value RedisQueueValue) GetTime() (t time.Time) {
	t = time.Unix(value.Time/1000, (value.Time%1000)*1000)
	return
}

func (value *RedisQueueValue) FromTime(t time.Time) {
	value.Time = t.Unix()*1000 + int64(t.Nanosecond()/1000000)
	return
}
