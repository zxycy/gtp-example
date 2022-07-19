package test

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	gtp2 "gtp-example/gtpv2_test/gtp"
	"net"
	"testing"
)

var conv2 *net.UDPConn

func init() {
	srvAddr, err := net.ResolveUDPAddr("udp", gtp2.GTPAddr)
	conv2, err = net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}
func TestV2EchoRequest(t *testing.T) {
	cbr := message.NewEchoRequest(5)
	er, _ := cbr.Marshal()
	fmt.Println(conv2)
	_, err := conv2.Write(er) // 发送数据
	if err != nil {
		fmt.Println("err: ", err)
	}
}
func TestCreateSessionRequest(t *testing.T) {
	cbr := message.NewCreateSessionRequest(5, 6)
	er, _ := cbr.Marshal()
	fmt.Println(conv2)
	_, err := conv2.Write(er) // 发送数据
	if err != nil {
		fmt.Println("err: ", err)
	}
}
