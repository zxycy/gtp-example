package gtp

import (

	"encoding/binary"
	"net"
	"github.com/skoef/ndp"
)

func BuildRAMessage(TEID uint32, QFI byte, UEIPv6 net.IP, Prefix net.IP) []byte {

	//构造RA消息
	icmp := &ndp.ICMPRouterAdvertisement{
		HopLimit:         0,
		ManagedAddress:   false,
		OtherStateful:    false,
		HomeAgent:        false,
		RouterPreference: ndp.RouterPreferenceLow,
		RouterLifeTime:   7200,
		ReachableTime:    0,
		RetransTimer:     0,
	}

	prefixInfomation := &ndp.ICMPOptionPrefixInformation{
		PrefixLength:      64,
		OnLink:            true,
		Auto:              true,
		ValidLifetime:     0xffffffff,
		PreferredLifetime: 0xffffffff,
		Prefix:            Prefix,
	}
	icmp.AddOption(prefixInfomation)

	//add MTU
	mtuOption := &ndp.ICMPOptionMTU{
		MTU: 1500,
	}
	icmp.AddOption(mtuOption)
	//
	//计算校验和
	body, err := icmp.Marshal()


	RADSrcAddr := "fe80::1" //源地址使用本机的IPv6地址
	//RAeDstAddr := "ff02::1"//广播
	err = ndp.Checksum(&body, net.ParseIP(RADSrcAddr), UEIPv6)

	if err != nil {

	}

	//构造GTP消息
	gtpMessage := []byte{0x34, 0xff} //固定的GTP头部

	//写入gtp payload length
	gtpPayloadLength := make([]byte, 2)
	binary.BigEndian.PutUint16(gtpPayloadLength, uint16(len(body)+40+8)) //GTP的payload长度 = RA消息的长度 + 40字节的IP头长度 //8为 extension header
	gtpMessage = append(gtpMessage, gtpPayloadLength...)

	//写入隧道标识TEID
	teid := make([]byte, 4)
	binary.BigEndian.PutUint32(teid, TEID)
	gtpMessage = append(gtpMessage, teid...)

	//Next extension header type: PDU Session container (0x85)
	htype := make([]byte, 4)
	htype = []byte{0x00, 0x00, 0x00, 0x85}
	//Extension header (PDU Session container)

	exheader := make([]byte, 4)
	exheader = []byte{0x01, 0x00, QFI, 0x00}

	gtpMessage = append(gtpMessage, htype...)
	gtpMessage = append(gtpMessage, exheader...)

	//组装ICMPv6消息的IP头部，固定40个字节
	ipHead := []byte{0x60, 0x00, 0x00, 0x00} //固定值
	//payload
	icmpPayloadLength := make([]byte, 2) //ICMPv6消息payload长度就是RA消息的长度，动态填写
	binary.BigEndian.PutUint16(icmpPayloadLength, uint16(len(body)))
	ipHead = append(ipHead, icmpPayloadLength...)
	//nexthead
	nextHead := []byte{0x3a} //ICMPv6标识
	ipHead = append(ipHead, nextHead...)
	//hoplimit
	hoplimit := []byte{0xff}
	ipHead = append(ipHead, hoplimit...)
	//源地址和目的地址
	sourIP := net.ParseIP(RADSrcAddr)

	//destIP := net.ParseIP(RAeDstAddr)//使用广播地址
	ipHead = append(ipHead, sourIP...)
	ipHead = append(ipHead, UEIPv6...)

	//组装UDP消息消息的payload: gtp头 + ip头 + RA消息

	gtpMessage = append(gtpMessage, ipHead...)
	gtpMessage = append(gtpMessage, body...)

	return gtpMessage
}

func BuildRS() []byte {
	rs:=[]byte{52,255,0,56,0,0,0,5,0,0,0, 133,1,0,4,0,
		96,0,0,0,0,8,58,255,254,128,0,0,0,0,0,0,0,0,0,0,0,0,0,3,255,2,0,0,0,0,0,0,0,0,0,0,0,0,0,2,
		133,0, 125,52,0,0,0,0}
	return rs
}

func BuildPing()[]byte  {
	ping:=[]byte{52,255,0,112,0,0,0,16,0,0,0, 133,1,0,4,0,
		96,7,13,0,0,64,58,63,32,1,4,112,130,135,6,23,216,107,141,7,36,111,93,61,36,7,174,128,2,0,16,1,0,0,0,0,0,0,0,32,
		128,0,148,226,146,235,0,0,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,57,57,32,98,
		111,116,116,108,101,115,32,111,102,32,98,101,101,114,32,111,110,32,116,104,101,32,119,97,108,108}
	return ping
}
func BuildReply() []byte {
	ping:=[]byte{52,255,0,112,0,0,0,16,0,0,0, 133,1,0,4,0,
		96,7,13,0,0,64,58,63,36,7,174,128,2,0,16,1,0,0,0,0,0,0,0,32,32,1,4,112,130,135,6,23,216,107,141,7,36,111,93,61,
		129,0,147,226,146,235,0,0,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,57,57,32,98,
		111,116,116,108,101,115,32,111,102,32,98,101,101,114,32,111,110,32,116,104,101,32,119,97,108,108}
	return ping
}