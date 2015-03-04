package main

import (
	"flag"
	"fmt"

	"encoding/binary"
	"github.com/fanngyuan/link"
	"math/rand"
	"time"
	"zerogame.info/taphero/app"
	"zerogame.info/taphero/conf"
	"zerogame.info/taphero/log"
	"zerogame.info/taphero/logic"
	"zerogame.info/taphero/net"
)

const (
	PONG = "pong"
	PING = "ping"
	MB   = 1073741824
)

func main() {
	mode := flag.String("mode", "dev", "run mode")
	// platform := flag.Uint64("platform", 1, "platform")
	server := flag.Uint64("server", 1, "server")
	process := flag.Uint64("process", 0, "process")
	locale := flag.String("locale", "chi", "locale[chi,eng...]")
	flag.Parse()

	app.InitResource(*mode, *server, *process, *locale)
	proto := link.PacketN(4, link.BigEndian)
	proto.MaxPacketSize = MB
	linkServer, err := link.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", "7788"), proto)
	if err != nil {
		panic(err)
	}

	linkServer.Handle(func(session *link.Session) {
		// log.Info("client", session.Conn().RemoteAddr().String(), "in")
		lastTime := time.Now()
		rand.Seed(lastTime.UnixNano())
		go func() {
			for {
				time.Sleep(getPingDur() * time.Second)
				data := make([]byte, 8)
				binary.BigEndian.PutUint32(data[0:4], 4)
				copy(data[4:], []byte(PONG))
				session.SendPacket(net.WritePacket(data))
				// session.SendPacket(net.WritePacket([]byte(PONG)))
				now := time.Now()
				if now.Sub(lastTime) > getSessionTimeout()*time.Second {
					app.CloseTcp(session, now.Add(-getSessionTimeout()*time.Second), "pang")
					break
				}
			}
		}()
		session.Handle(func(buffer *link.InBuffer) {
			msg := buffer.Data
			//log.Debug("client", session.Conn().RemoteAddr().String(), "say:", string(msg))
			// log.Debug(session.Conn().RemoteAddr().String())
			lastTime = time.Now()
			if len(msg) == 4 && string(msg) == PING {
				log.Debug("get_ping", session.Id(), session.Conn().RemoteAddr().String())
				log.Debug("get ping", session.Id(), lastTime)
				// do nothing
			} else {
				logic.HandleMessage(msg, session, lastTime)
			}
		})

		app.CloseTcp(session, time.Now(), "tcpclosed")
	})

}
func getPingDur() time.Duration {
	if conf.IsModePro() {
		return 60
	}
	return 5

}
func getSessionTimeout() time.Duration {
	if conf.IsModePro() {
		return 60 * 2
	}
	return 30

}
