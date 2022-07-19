package pducontainer

import (
	"encoding/binary"
	"fmt"
	"time"
)

const PDUType_UPPDUSESSIONINFORMATION = 1

type ULPDUSESSIONINFORMATION struct {
	PDUType    uint8 //0x11110000
	QMP        bool  //0x00001000
	DLDelayInd bool  //0x00000100
	ULDelayInd bool  //0x00000010
	SNP        bool  //0x00000001

	N3N9DelayInd      bool  //0x10000000
	NewIEFlag         bool  //0x01000000
	QoSFlowIdentifier uint8 //0x00111111

	DLReceivedTimeStamp        time.Time
	ULSendingTimeStamp         time.Time
	DLSendingTimeStampRepeated time.Time
	DLDelayResult              uint32
	ULDelayResult              uint32
	ULQFISequenceNumber        []byte
	N3N9DelayResult            uint32
	NewIEFlags                 []bool
}

func NewUlPduSessionInfo(qfi uint8) ULPDUSESSIONINFORMATION {
	return ULPDUSESSIONINFORMATION{
		PDUType:           PDUType_UPPDUSESSIONINFORMATION,
		QoSFlowIdentifier: qfi,
	}
}

func (u *ULPDUSESSIONINFORMATION) MarshalBinary() (data []byte, err error) {
	var idx uint16 = 0

	// Octet 1
	tmpUint8 := u.PDUType<<4 |
		btou(u.QMP)<<3 |
		btou(u.DLDelayInd)<<2 |
		btou(u.ULDelayInd)<<1 |
		btou(u.SNP)
	data = append([]byte(""), byte(tmpUint8))
	idx = idx + 1

	// Octet 2
	tmpUint8 = btou(u.N3N9DelayInd)<<7 |
		btou(u.NewIEFlag)<<6 |
		u.QoSFlowIdentifier
	data = append(data, byte(tmpUint8))
	idx = idx + 1

	// Octet +8*3
	if u.QMP {
		duration := u.DLReceivedTimeStamp.Sub(BASE_DATE).Seconds()
		data = append(data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(data[idx:], uint64(duration))
		idx = idx + 8

		duration = u.ULSendingTimeStamp.Sub(BASE_DATE).Seconds()
		data = append(data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(data[idx:], uint64(duration))
		idx = idx + 8

		duration = u.DLSendingTimeStampRepeated.Sub(BASE_DATE).Seconds()
		data = append(data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(data[idx:], uint64(duration))
		idx = idx + 8
	}

	// Octet +4
	if u.DLDelayInd {
		data = append(data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(data[idx:], u.DLDelayResult)
		idx = idx + 4
	}

	// Octet +4
	if u.ULDelayInd {
		data = append(data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(data[idx:], u.ULDelayResult)
		idx = idx + 4
	}

	// Octet +3
	if u.SNP {
		if len(u.ULQFISequenceNumber) != 3 {
			return []byte(""), fmt.Errorf("Length of DLQFISequenceNumber data shall be 3 ")
		}
		data = append(data, u.ULQFISequenceNumber...)
		idx = idx + 3
	}

	// Octet +4
	if u.N3N9DelayInd {
		data = append(data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(data[idx:], u.N3N9DelayResult)
		idx = idx + 4
	}

	if u.NewIEFlag {

	}
	return data, nil
}
func (u *ULPDUSESSIONINFORMATION) UnmarshalBinary(data []byte) error {
	length := uint16(len(data))

	var idx uint16 = 0
	// Octet 1
	if length < idx+1 {
		return fmt.Errorf("Inadequate TLV length: %d", length)
	}
	u.PDUType = data[idx] >> 4
	u.QMP = utob(uint8(data[idx]) & BitMask4)
	u.DLDelayInd = utob(uint8(data[idx]) & BitMask3)
	u.ULDelayInd = utob(uint8(data[idx]) & BitMask2)
	u.SNP = utob(uint8(data[idx]))
	idx = idx + 1

	// Octet 2
	u.N3N9DelayInd = utob(uint8(data[idx]) & BitMask8)
	u.NewIEFlag = utob(uint8(data[idx]) & BitMask7)
	u.QoSFlowIdentifier = data[idx] & Mask3
	idx = idx + 1

	if u.QMP {
		if length < idx+24 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		duration := binary.BigEndian.Uint64(data[idx : idx+8])
		u.DLReceivedTimeStamp = BASE_DATE.Add(time.Duration(duration) * time.Second)
		idx = idx + 8

		duration = binary.BigEndian.Uint64(data[idx : idx+8])
		u.ULSendingTimeStamp = BASE_DATE.Add(time.Duration(duration) * time.Second)
		idx = idx + 8

		duration = binary.BigEndian.Uint64(data[idx : idx+8])
		u.DLSendingTimeStampRepeated = BASE_DATE.Add(time.Duration(duration) * time.Second)
		idx = idx + 8
	}

	if u.DLDelayInd {
		if length < idx+4 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		u.DLDelayResult = binary.BigEndian.Uint32(data[idx : idx+4])
		idx = idx + 4
	}

	if u.ULDelayInd {
		if length < idx+4 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		u.ULDelayResult = binary.BigEndian.Uint32(data[idx : idx+4])
		idx = idx + 4
	}

	if u.SNP {
		if length < idx+3 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		u.ULQFISequenceNumber = data[idx : idx+3]
		idx = idx + 3
	}

	if u.N3N9DelayInd {
		if length < idx+4 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		u.N3N9DelayResult = binary.BigEndian.Uint32(data[idx : idx+4])
		idx = idx + 4
	}

	if length != idx {
		return fmt.Errorf("Inadequate TLV length: %d ", length)
	}
	return nil
}
