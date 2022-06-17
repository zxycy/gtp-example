package pducontainer

import (
	"fmt"
	"testing"
	"time"
)

func TestDLPDUSESSIONINFORMATION(t *testing.T) {
	dlpdu := DLPDUSESSIONINFORMATION{
		PDUType:                0,
		QMP:                    true,
		SNP:                    true,
		MSNP:                   true,
		PPP:                    true,
		RQI:                    true,
		QoSFlowIdentifier:      66,
		PPI:                    1,
		DLSendingTimeStamp:     time.Now(),
		DLQFISequenceNumber:    []byte{8, 8, 8},
		DLMBSQFISequenceNumber: 3,
	}
	dlbody, err := dlpdu.MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("binary:", dlbody)

	var dlpdu2 DLPDUSESSIONINFORMATION
	dlBody := []byte{14, 194, 32, 0, 0, 0, 0, 230, 83, 224, 196, 8, 8, 8, 0, 0, 0, 3}
	if err = dlpdu2.UnmarshalBinary(dlBody); err != nil {
		fmt.Println(err)
	}
	fmt.Println("struct:", dlpdu2)

	dl := NewDlPduSessionInfo(5)
	fmt.Println("struct2:", dl)
}

func TestULPDUSESSIONINFORMATION(t *testing.T) {
	uppdu := ULPDUSESSIONINFORMATION{
		PDUType:                    1,
		QMP:                        true,
		DLDelayInd:                 true,
		ULDelayInd:                 true,
		SNP:                        true,
		N3N9DelayInd:               true,
		QoSFlowIdentifier:          10,
		DLReceivedTimeStamp:        time.Now(),
		ULSendingTimeStamp:         time.Now(),
		DLSendingTimeStampRepeated: time.Now(),
		DLDelayResult:              0,
		ULDelayResult:              0,
		ULQFISequenceNumber:        []byte{6, 6, 6},
		N3N9DelayResult:            525,
	}
	ulbody, err := uppdu.MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("binary:", ulbody)

	var ulpdu2 ULPDUSESSIONINFORMATION
	llBody := []byte{31, 138, 0, 0, 0, 0, 230, 83, 245, 96, 0, 0, 0, 0, 230, 83, 247, 216,
		0, 0, 0, 0, 230, 84, 3, 105, 0, 0, 0, 0, 0, 0, 0, 0, 6, 6, 6, 0, 0, 2, 13}
	if err = ulpdu2.UnmarshalBinary(llBody); err != nil {
		fmt.Println(err)
	}
	fmt.Println("struct:", ulpdu2)

	ul := NewUlPduSessionInfo(5)
	fmt.Println("struct2:", ul)
}
