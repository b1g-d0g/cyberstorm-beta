package core

import (
	"encoding/binary"
	"net"
)

// GetUDPHeaderBytes builds the UDP packet header
func GetUDPHeaderBytes(srcIP, dstIP net.IP, dstPort uint16) ([]byte, error) {
	srcPort := uint16(12345) // Random source port

	udpLength := uint16(8) // UDP header length is always 8 bytes

	b := make([]byte, 8)
	binary.BigEndian.PutUint16(b[0:2], srcPort) // Source port
	binary.BigEndian.PutUint16(b[2:4], dstPort) // Destination port
	binary.BigEndian.PutUint16(b[4:6], udpLength)
	binary.BigEndian.PutUint16(b[6:8], 0) // Checksum

	return b, nil
}
