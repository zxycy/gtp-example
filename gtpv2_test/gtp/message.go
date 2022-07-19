package gtp

import (
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"net"
)

type GTPv2Message struct {
	RemoteAddr *net.UDPAddr
	GtpMessage message.Message
	BinaryMsg  []byte
}

func NewGTPv2Message(addr *net.UDPAddr, gtpMsg message.Message, bMsg []byte) GTPv2Message {
	return GTPv2Message{
		RemoteAddr: addr,
		GtpMessage: gtpMsg,
		BinaryMsg:  bMsg,
	}
}
