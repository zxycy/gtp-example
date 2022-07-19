package test

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	gtp1 "gtp-example/gtpv1_test/gtp"
	"net"
	"testing"
)

var con *net.UDPConn

func init() {
	srvAddr, err := net.ResolveUDPAddr("udp", gtp1.GTPAddr)
	con, err = net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}
func TestPing(t *testing.T) {
	ping := gtp1.BuildPing()
	fmt.Println(con)
	_, err := con.Write(ping)
	if err != nil {
		fmt.Println(err)
	}

}
func TestRouterSolicitation(t *testing.T) {
	rs := gtp1.BuildRS()
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
