package gtp

import (
	"fmt"
	"github.com/mdlayher/ndp"
	"github.com/wmnsk/go-gtp/gtpv1/message"
	"gtp-example/pducontainer"
	"net"
	"net/netip"
)

func BuildRAMessage(TEID uint32, QFI byte, UEIPv6 net.IP, Prefix netip.Prefix) []byte {
	//构造RA消息
	fmt.Println("Build RouterAdvertisement Message")
	icmp := &ndp.RouterAdvertisement{
		CurrentHopLimit:           0,
		ManagedConfiguration:      false,
		OtherConfiguration:        false,
		MobileIPv6HomeAgent:       false,
		RouterSelectionPreference: ndp.Low,
		RouterLifetime:            7200,
		ReachableTime:             0,
		RetransmitTimer:           0,
		Options:                   nil,
	}

	//add PrefixInfomation
	prefixInfomation := &ndp.PrefixInformation{
		PrefixLength:                   64,
		OnLink:                         true,
		AutonomousAddressConfiguration: true,
		ValidLifetime:                  0xffffffff,
		PreferredLifetime:              0xffffffff,
		Prefix:                         Prefix.Addr(),
	}
	icmp.Options = append(icmp.Options, prefixInfomation)

	//add MTU
	mtuOption := &ndp.MTU{
		MTU: 1500,
	}
	icmp.Options = append(icmp.Options, mtuOption)

	cfgIP := "fe80::1" //本机ip
	SrcAddr, _ := netip.ParseAddr(cfgIP)
	DstAddr, _ := netip.ParseAddr(UEIPv6.String())
	body, err := ndp.MarshalMessageChecksum(icmp, SrcAddr, DstAddr)
	if err != nil {
		fmt.Println(err)
	}

	//IPv6 header
	v6header, _ := pducontainer.NewICMPv6Header(len(body), net.ParseIP(cfgIP), UEIPv6)

	//gtp message
	var tempBody []byte
	tempBody = append(tempBody, v6header...)
	tempBody = append(tempBody, body...)
	gtpMessage := message.NewTPDU(TEID, tempBody)

	//add gtp extension header
	dlPdu := pducontainer.NewDlPduSessionInfo(QFI)
	container, err := dlPdu.MarshalBinary()
	exheader := message.NewExtensionHeader(message.ExtHeaderTypePDUSessionContainer, container, 0)
	gtpMessage.AddExtensionHeaders(exheader)

	gtpBody, _ := gtpMessage.Marshal()

	return gtpBody
}

func BuildRS() []byte {
	fmt.Println("Build rs Message")
	icmp := ndp.RouterSolicitation{Options: nil}
	SrcAddr, _ := netip.ParseAddr("fe80::3")
	DstAddr, _ := netip.ParseAddr("ff02::2")
	body, _ := ndp.MarshalMessageChecksum(&icmp, SrcAddr, DstAddr)

	v6header, _ := pducontainer.NewICMPv6Header(len(body), net.ParseIP("fe80::3"), net.ParseIP("ff02::2"))

	var tempBody []byte
	tempBody = append(tempBody, v6header...)
	tempBody = append(tempBody, body...)
	gtpMessage := message.NewTPDU(3, tempBody)

	//add gtp extension header
	dlPdu := pducontainer.NewDlPduSessionInfo(3)
	container, _ := dlPdu.MarshalBinary()
	exheader := message.NewExtensionHeader(message.ExtHeaderTypePDUSessionContainer, container, 0)
	gtpMessage.AddExtensionHeaders(exheader)

	gtpBody, _ := gtpMessage.Marshal()

	return gtpBody

}

func BuildPing() []byte {
	ping := []byte{52, 255, 0, 112, 0, 0, 0, 16, 0, 0, 0, 133, 1, 0, 4, 0,
		96, 7, 13, 0, 0, 64, 58, 63, 32, 1, 4, 112, 130, 135, 6, 23, 216, 107, 141, 7, 36, 111, 93, 61, 36, 7, 174, 128, 2, 0, 16, 1, 0, 0, 0, 0, 0, 0, 0, 32,
		128, 0, 148, 226, 146, 235, 0, 0, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 57, 57, 32, 98,
		111, 116, 116, 108, 101, 115, 32, 111, 102, 32, 98, 101, 101, 114, 32, 111, 110, 32, 116, 104, 101, 32, 119, 97, 108, 108}
	return ping
}
func BuildReply() []byte {
	ping := []byte{52, 255, 0, 112, 0, 0, 0, 16, 0, 0, 0, 133, 1, 0, 4, 0,
		96, 7, 13, 0, 0, 64, 58, 63, 36, 7, 174, 128, 2, 0, 16, 1, 0, 0, 0, 0, 0, 0, 0, 32, 32, 1, 4, 112, 130, 135, 6, 23, 216, 107, 141, 7, 36, 111, 93, 61,
		129, 0, 147, 226, 146, 235, 0, 0, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 57, 57, 32, 98,
		111, 116, 116, 108, 101, 115, 32, 111, 102, 32, 98, 101, 101, 114, 32, 111, 110, 32, 116, 104, 101, 32, 119, 97, 108, 108}
	return ping
}
