package resource

import (
	key "github.com/0studio/storage_key"
	"zerogame.info/taphero/entity"
)

var sessionUserMap map[uint64]key.KeyUint64 // key session.Id,value user.Uin
var sessions entity.GameSessionMap          // key= user.Uin,value-session

func Init() {
	sessionUserMap = make(map[uint64]key.KeyUint64)
	sessions = make(map[key.KeyUint64]entity.GameSession)
}
func PutUin(tcpSessionId uint64, userUin key.KeyUint64) {
	sessionUserMap[tcpSessionId] = userUin
}
func GetUin(tcpSessionId uint64) (userUin key.KeyUint64, ok bool) {
	userUin, ok = sessionUserMap[tcpSessionId]
	return
}
func DeleteUin(tcpSessionId uint64) {
	delete(sessionUserMap, tcpSessionId)
}
func DelGameSession(userUin key.KeyUint64) {
	gameSession, ok := GetGameSession(userUin)
	if ok {
		gameSession.LinkSession.Conn().Close()
		gameSession.LinkSession.Close("del")
		DeleteUin(gameSession.LinkSession.Id())
		delete(sessions, userUin)
	}

}
func GetGameSession(userUin key.KeyUint64) (gameSession entity.GameSession, ok bool) {
	gameSession, ok = sessions[userUin]
	return
}
func PutGameSession(userUin key.KeyUint64, gameSession entity.GameSession) {
	sessions[userUin] = gameSession
}
func GetGameSessionByTcpSessionId(tcpSessionId uint64) (gameSession entity.GameSession, ok bool) {
	userUin, ok := GetUin(tcpSessionId)
	if !ok {
		return
	}
	gameSession, ok = GetGameSession(userUin)
	return

}
