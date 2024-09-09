package core

import (
	"encoding/binary"
	//"fmt"
	"math/rand"
	//"net"
	//"syscall"
)

// Costanti per i messaggi ICMP
const (
	ICMPEchoRequest = 8
	ICMPEchoReply   = 0
	ICMPHeaderLen   = 8
)

// GetICMPHeaderBytes crea l'header di un pacchetto ICMP Echo Request
func GetICMPHeaderBytes() ([]byte, error) {
	icmpHeader := make([]byte, ICMPHeaderLen)
	icmpHeader[0] = ICMPEchoRequest // Tipo: Echo Request
	icmpHeader[1] = 0               // Codice: 0 per Echo Request

	// Identifier e Sequence Number casuali
	identifier := rand.Intn(1<<16 - 1)
	sequenceNumber := rand.Intn(1<<16 - 1)

	// Impostiamo identifier e sequence number
	binary.BigEndian.PutUint16(icmpHeader[4:6], uint16(identifier))
	binary.BigEndian.PutUint16(icmpHeader[6:8], uint16(sequenceNumber))

	// Calcoliamo il checksum
	checksum := Checksum(icmpHeader)
	binary.BigEndian.PutUint16(icmpHeader[2:4], checksum) // Inseriamo il checksum nell'header

	return icmpHeader, nil
}

// Funzione per calcolare il checksum ICMP
func Checksum(data []byte) uint16 {
	var csum uint32
	for i := 0; i < len(data)-1; i += 2 {
		csum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[len(data)-1]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}
