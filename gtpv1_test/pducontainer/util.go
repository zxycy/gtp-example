package pducontainer

import (
	"encoding/binary"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

const (
	Mask8 = 1<<8 - 1
	Mask7 = 1<<7 - 1
	Mask6 = 1<<6 - 1
	Mask5 = 1<<5 - 1
	Mask4 = 1<<4 - 1
	Mask3 = 1<<3 - 1
	Mask2 = 1<<2 - 1
	Mask1 = 1<<1 - 1
)

const (
	BitMask8 = 1 << 7
	BitMask7 = 1 << 6
	BitMask6 = 1 << 5
	BitMask5 = 1 << 4
	BitMask4 = 1 << 3
	BitMask3 = 1 << 2
	BitMask2 = 1 << 1
	BitMask1 = 1
)

var BASE_DATE time.Time = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

func btou(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func utob(u uint8) bool {
	return u != 0
}

func ipv6HeaderMarshal(h *ipv6.Header) ([]byte, error) {
	b := make([]byte, 40)
	b[0] |= byte(h.Version) << 4
	b[0] |= byte(h.TrafficClass) >> 4
	b[1] |= byte(h.TrafficClass) << 4
	b[1] |= byte(h.FlowLabel >> 16)
	b[2] = byte(h.FlowLabel >> 8)
	b[3] = byte(h.FlowLabel)
	binary.BigEndian.PutUint16(b[4:6], uint16(h.PayloadLen))
	b[6] = byte(h.NextHeader)
	b[7] = byte(h.HopLimit)
	copy(b[8:24], h.Src)
	copy(b[24:40], h.Dst)
	return b, nil
}

func NewICMPv6Header(length int, srcIP, dstIP net.IP) ([]byte, error) {
	header := ipv6.Header{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		PayloadLen:   length,
		NextHeader:   58, //ICMPV6
		HopLimit:     0xff,
		Src:          srcIP,
		Dst:          dstIP,
	}
	return ipv6HeaderMarshal(&header)
}
