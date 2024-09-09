package core

import (
	"encoding/binary"
)

// GetICMPHeaderBytes builds the ICMP packet header
func GetICMPHeaderBytes() ([]byte, error) {
	b := make([]byte, 8)  // ICMP header size is 8 bytes
	b[0] = 8              // Type (8 = Echo request)
	b[1] = 0              // Code
	binary.BigEndian.PutUint16(b[6:8], 0) // Identifier and Sequence

	return b, nil
}
