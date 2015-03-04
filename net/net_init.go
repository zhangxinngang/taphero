package net

import (
	"github.com/fanngyuan/link"
)

func WritePacket(message []byte) link.Packet {
	buffer := new(link.OutBuffer)
	buffer.WriteBytes(message)
	return link.Packet{buffer}
}
