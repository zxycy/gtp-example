package main

import (
	"fmt"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"gtp-example/gtp"
	"gtp-example/server"
	"net"
	"time"
)

func main()  {

	go server.GTPServer()
	//go server.GTPClient()

	srvAddr, _ := net.ResolveUDPAddr("udp", server.ServerAddr)


	con, err:= net.DialUDP("udp", nil, srvAddr)
	if err!= nil {
		fmt.Println("err: ", err)
		return
	}
	defer con.Close()

	rs:= gtp.BuildRS()
	echo:=message.NewEchoRequest(5)
	er,_:=echo.Marshal()
	_, err = con.Write(rs) // 发送数据
	_, err = con.Write(er) // 发送数据
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	time.Sleep(100000000000000000)



}

//func Clinet()  {
//
//}