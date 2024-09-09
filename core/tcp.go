package core

import (
	"encoding/binary"
	"net"
)

// GetTCPSYNHeaderBytes builds the TCP SYN packet header
func GetTCPSYNHeaderBytes(srcIP, dstIP net.IP, dstPort uint16) ([]byte, error) {
	srcPort := uint16(12345) // Random source port for now

	seq := uint32(0) // Random initial sequence number
	offsetAndFlags := []byte{
		byte(5 << 4), // Data offset (header length / 4)
		0x02,        // SYN flag set
	}

	b := make([]byte, 20) // TCP header size is 20 bytes
	binary.BigEndian.PutUint16(b[0:2], srcPort)        // Source port
	binary.BigEndian.PutUint16(b[2:4], dstPort)        // Destination port
	binary.BigEndian.PutUint32(b[4:8], seq)            // Sequence number
	binary.BigEndian.PutUint32(b[8:12], 0)             // Acknowledgement number
	copy(b[12:14], offsetAndFlags)                     // Offset, reserved and flags
	binary.BigEndian.PutUint16(b[14:16], 65535)        // Window size
	binary.BigEndian.PutUint16(b[18:20], 0)            // Urgent pointer

	checksum, err := tcpChecksum(srcIP, dstIP, b)
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint16(b[16:18], checksum) // Checksum
	return b, nil
}

func tcpChecksum(srcIP, dstIP net.IP, data []byte) (uint16, error) {
	// Calculate TCP checksum here (pseudocode)
	// ...
	return 0, nil
}
