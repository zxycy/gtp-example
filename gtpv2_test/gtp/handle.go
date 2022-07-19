package gtp

import (
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func HandleCreateSessionRequest(msg GTPv2Message) {
	//c, _ := m2.ParseCreateSessionRequest(msg.BinaryMsg)

	res := message.NewCreateSessionResponse(5, 6)
	m, _ := res.Marshal()

	Sendgtp(m, msg.RemoteAddr)
}

func HandleEchoRequest(msg GTPv2Message) {
	er, _ := message.ParseEchoRequest(msg.BinaryMsg)

	res := message.NewEchoResponse(er.SequenceNumber)

	m, _ := res.Marshal()

	Sendgtp(m, msg.RemoteAddr)
}