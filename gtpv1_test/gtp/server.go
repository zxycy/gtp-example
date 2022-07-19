package gtp

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"net"
)

const GTPAddr = "127.0.0.2:2152"

func Run() {
	go GTPServer()
}

func GTPServer() {

	var udpaddr *net.UDPAddr

	udpaddr, _ = net.ResolveUDPAddr("udp4", GTPAddr)

	udpconn, _ := net.ListenUDP("udp", udpaddr)
	go func() {
		for {

			var buf = make([]byte, 1024)
			_, addr, err := udpconn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			switch buf[1] {
			case message.MsgTypeEchoRequest:
				echoreq, _ := message.ParseEchoRequest(buf)
				HandlerEchoRequest(echoreq, addr, udpconn)
			case message.MsgTypeTPDU:
				tpduMsg, _ := message.ParseTPDU(buf)
				HandlerTPDU(tpduMsg, addr, udpconn)
			default:

			}
		}
	}()

}
