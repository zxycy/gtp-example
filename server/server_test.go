package server

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"gtp-example/gtp"
	"net"
	"testing"
)

var con *net.UDPConn

func init() {
	srvAddr, err := net.ResolveUDPAddr("udp", ServerAddr)

	con, err = net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}
func TestPing(t *testing.T) {
	ping := gtp.BuildPing()
	_, err := con.Write(ping)
	if err != nil {
		fmt.Println(err)
	}

}
func TestRouterSolicitation(t *testing.T) {
	rs := gtp.BuildRS()
	_, err := con.Write(rs) // 发送数据
	if err != nil {
		fmt.Println(err)
	}

}

func TestEchoRequest(t *testing.T) {
	echo := message.NewEchoRequest(5)
	er, _ := echo.Marshal()
	_, err := con.Write(er) // 发送数据
	if err != nil {
		fmt.Println("err: ", err)
	}
}
