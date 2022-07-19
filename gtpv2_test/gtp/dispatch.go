package gtp

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func Dispatch(msg GTPv2Message) {
	fmt.Println(msg.GtpMessage.MessageType())
	switch msg.GtpMessage.MessageType() {
	case message.MsgTypeEchoRequest:
		fmt.Printf("ffff")
		HandleEchoRequest(msg)
	case message.MsgTypeEchoResponse:
	case message.MsgTypeCreateSessionRequest:
		HandleCreateSessionRequest(msg)
	case message.MsgTypeCreateBearerResponse:
	case message.MsgTypeDeleteSessionRequest:
	case message.MsgTypeDeleteSessionResponse:
	case message.MsgTypeChangeNotificationRequest:
	case message.MsgTypeChangeNotificationResponse:
	case message.MsgTypeRemoteUEReportNotification:
	case message.MsgTypeRemoteUEReportAcknowledge:

	default:
		fmt.Printf("Unknown GTP message type: %d", msg.GtpMessage.MessageType())
		return
	}
}
