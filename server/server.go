package server

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"gtp-example/gtp"
	"net"
)

const (
		ServerAddr  = "127.0.0.10:2152"
		ClientAddr  = "127.0.0.1:2152"
	)

func GTPServer() {



	var udpaddr *net.UDPAddr

	udpaddr, _ = net.ResolveUDPAddr("udp4", ServerAddr)

	udpconn, _ := net.ListenUDP("udp", udpaddr)

	for {

		var buf = make([]byte, 1024)
		_, addr, err := udpconn.ReadFromUDP(buf)
		if err!=nil{
			fmt.Println("err:",err)
			return
		}
		switch buf[1] {
		case message.MsgTypeEchoRequest:
			echoreq, _ := message.ParseEchoRequest(buf)
			gtp.HandlerEchoRequest(echoreq, addr, udpconn)
		case message.MsgTypeTPDU:
			tpduMsg, _ := message.ParseTPDU(buf)
			gtp.HandlerTPDU(tpduMsg, addr, udpconn)
		default:

		}
	}

}

func GTPClient()  {

	var udpaddr *net.UDPAddr

	udpaddr, _ = net.ResolveUDPAddr("udp4", ClientAddr)

	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err!=nil{
		fmt.Println(udpconn)
	}
	
}
