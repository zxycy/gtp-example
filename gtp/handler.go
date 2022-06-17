package gtp

import (
	"fmt"
	"github.com/mdlayher/ndp"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"golang.org/x/net/ipv6"
	"net"
	"net/netip"
)

func HandlerTPDU(gtpMsg *message.TPDU, addr *net.UDPAddr, udpconn *net.UDPConn) {
	var QFI byte
	for _, exheader := range gtpMsg.ExtensionHeaders {
		if exheader.Type == message.ExtHeaderTypePDUSessionContainer {
			QFI = exheader.Content[1] & 0x3f //这个字节的前六位为QFI
		}
	}

	v6Msg, _ := ipv6.ParseHeader(gtpMsg.Payload)
	ueIP := v6Msg.Src

	if v6Msg.NextHeader == 58 { //ICMPV6
		icmpv6Buf := gtpMsg.Payload[40:]
		msg, err := ndp.ParseMessage(icmpv6Buf)

		// todo
		if err != nil && icmpv6Buf[0] == byte(ipv6.ICMPTypeEchoRequest) {
			fmt.Println("received unexpected message but handle, err: ", err)
			pingReply := BuildReply()
			if _, err := udpconn.WriteToUDP(pingReply, addr); err != nil {
				fmt.Println("err:", err)
			}
			return
		}
		switch msg.Type() {
		case ipv6.ICMPTypeRouterSolicitation:
			prefix, _ := netip.ParsePrefix("1::/64")
			RA := BuildRAMessage(5, QFI, ueIP, prefix)
			//addr.Port = 2152
			if _, err := udpconn.WriteToUDP(RA, addr); err != nil {
				fmt.Println("err:", err)
			}
		case ipv6.ICMPTypeRouterAdvertisement:
			fmt.Println("not implemented")
		default:
			fmt.Println("not implemented")

		}

	}
}

func HandlerEchoRequest(msg *message.EchoRequest, addr *net.UDPAddr, udpconn *net.UDPConn) {

	rspbody := message.NewEchoResponse(msg.SequenceNumber, nil)
	res, _ := rspbody.Marshal()
	if _, err := udpconn.WriteToUDP(res, addr); err != nil {
		fmt.Println("err:", err)
	}
}
