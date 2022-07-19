package pducontainer

import (
	"encoding/binary"
	"fmt"
	"time"
)

const PDUType_DLPDUSESSIONINFORMATION = 0

type DLPDUSESSIONINFORMATION struct {
	PDUType uint8 //0x11110000
	QMP     bool  //0x00001000
	SNP     bool  //0x00000100
	MSNP    bool  //0x00000010
	//Spare       //0x00000001
	PPP               bool  //0x10000000
	RQI               bool  //0x01000000
	QoSFlowIdentifier uint8 //0x00111111

	PPI uint8 //0x11100000
	//Spare		//0x00011111
	DLSendingTimeStamp     time.Time
	DLQFISequenceNumber    []byte //3
	DLMBSQFISequenceNumber uint32
}

func NewDlPduSessionInfo(qfi uint8) DLPDUSESSIONINFORMATION {
	return DLPDUSESSIONINFORMATION{
		PDUType:           PDUType_DLPDUSESSIONINFORMATION,
		QoSFlowIdentifier: qfi,
	}
}
func (d *DLPDUSESSIONINFORMATION) MarshalBinary() (data []byte, err error) {
	var idx uint16 = 0

	// Octet 1
	tmpUint8 := d.PDUType<<4 |
		btou(d.QMP)<<3 |
		btou(d.SNP)<<2 |
		btou(d.MSNP)<<1
	data = append([]byte(""), byte(tmpUint8))
	idx = idx + 1

	// Octet 2
	tmpUint8 = btou(d.PPP)<<7 |
		btou(d.RQI)<<6 |
		d.QoSFlowIdentifier
	data = append(data, byte(tmpUint8))
	idx = idx + 1

	// Octet 3
	if d.PPP {
		tmpUint8 = d.PPI << 5
		data = append(data, byte(tmpUint8))
		idx = idx + 1
	}
	// Octet +8
	if d.QMP {
		duration := d.DLSendingTimeStamp.Sub(BASE_DATE).Seconds()
		data = append(data, make([]byte, 8)...)
		binary.BigEndian.PutUint64(data[idx:], uint64(duration))
		idx = idx + 8
	}

	// Octet +3
	if d.SNP {
		if len(d.DLQFISequenceNumber) != 3 {
			return []byte(""), fmt.Errorf("Length of DLQFISequenceNumber data shall be 3 ")
		}
		data = append(data, d.DLQFISequenceNumber...)
		idx = idx + 3
	}

	if d.MSNP {
		data = append(data, make([]byte, 4)...)
		binary.BigEndian.PutUint32(data[idx:], d.DLMBSQFISequenceNumber)
		idx = idx + 4
	}
	return data, nil
}
func (d *DLPDUSESSIONINFORMATION) UnmarshalBinary(data []byte) error {
	length := uint16(len(data))

	var idx uint16 = 0
	// Octet 1
	if length < idx+1 {
		return fmt.Errorf("Inadequate TLV length: %d", length)
	}
	d.PDUType = data[idx] >> 4
	d.QMP = utob(uint8(data[idx]) & BitMask4)
	d.SNP = utob(uint8(data[idx]) & BitMask3)
	d.MSNP = utob(uint8(data[idx]) & BitMask2)
	idx = idx + 1

	// Octet 2
	d.PPP = utob(uint8(data[idx]) & BitMask7)
	d.RQI = utob(uint8(data[idx]) & BitMask6)
	d.QoSFlowIdentifier = data[idx] & Mask7
	idx = idx + 1

	if d.PPP {
		d.PPI = data[idx] >> 5
		idx = idx + 1
	}
	if d.QMP {
		if length < idx+8 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		duration := binary.BigEndian.Uint64(data[idx : idx+8])
		d.DLSendingTimeStamp = BASE_DATE.Add(time.Duration(duration) * time.Second)
		idx = idx + 8

	}

	if d.SNP {
		if length < idx+3 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		d.DLQFISequenceNumber = data[idx : idx+3]
		idx = idx + 3
	}
	if d.MSNP {
		if length < idx+4 {
			return fmt.Errorf("Inadequate TLV length: %d ", length)
		}
		d.DLMBSQFISequenceNumber = binary.BigEndian.Uint32(data[idx : idx+4])
		idx = idx + 4
	}

	if length != idx {
		return fmt.Errorf("Inadequate TLV length: %d ", length)
	}
	return nil
}
