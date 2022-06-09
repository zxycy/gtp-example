package gtp

import (
	"fmt"
	"github.com/skoef/ndp"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"golang.org/x/net/ipv6"
	"net"
)



func HandlerTPDU(msg *message.TPDU, addr *net.UDPAddr, udpconn *net.UDPConn) {
	var QFI byte
	for _, exheader := range msg.ExtensionHeaders {
		if exheader.Type == message.ExtHeaderTypePDUSessionContainer {
			QFI = exheader.Content[1] & 0x3f //这个字节的前六位为QFI
		}
	}

	v6Msg, _ := ipv6.ParseHeader(msg.Payload)
	ueIP := v6Msg.Src

	if v6Msg.NextHeader == 58 { //ICMPV6
		icmpv6Buf := msg.Payload[40:]
		msg, _ := ndp.ParseMessage(icmpv6Buf)
		if msg.Type() == ipv6.ICMPTypeRouterSolicitation {
			RA := BuildRAMessage(5, QFI, ueIP, net.ParseIP("fe80::1"))
			//addr.Port = 2152
			if _, err := udpconn.WriteToUDP(RA, addr); err != nil {
				fmt.Println("err:",err)
			}
		}
		if msg.Type() ==ipv6.ICMPTypeEchoRequest{
			pingReply:= BuildReply()
			if _, err := udpconn.WriteToUDP(pingReply, addr); err != nil {
				fmt.Println("err:",err)
			}
		}
	}
}


func HandlerEchoRequest(msg *message.EchoRequest, addr *net.UDPAddr, udpconn *net.UDPConn) {

	rspbody := message.NewEchoResponse(msg.SequenceNumber, nil)
	res, _ := rspbody.Marshal()
	if _, err := udpconn.WriteToUDP(res, addr); err != nil {
		fmt.Println("err:",err)
	}
}