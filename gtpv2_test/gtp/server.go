package gtp

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"net"
)

var GTPv2Con *net.UDPConn

const GTPAddr = "127.0.0.2:2123"

func Run(Dispatch func(GTPv2Message)) {
	var udpaddr *net.UDPAddr
	udpaddr, _ = net.ResolveUDPAddr("udp4", GTPAddr)

	GTPv2Con, _ = net.ListenUDP("udp", udpaddr)

	go func() {
		for {
			var buf = make([]byte, 1024)
			_, addr, err := GTPv2Con.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			msg, _ := message.Parse(buf)
			gtpmsg := NewGTPv2Message(addr, msg, buf)
			go Dispatch(gtpmsg)
		}
	}()

}

func Sendgtp(msg []byte, addr *net.UDPAddr) {
	_, err := GTPv2Con.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Println("err:", err)
	}

}
