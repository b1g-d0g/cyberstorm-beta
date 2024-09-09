package core

import (
	"encoding/binary"
	"fmt"
	"net"
)

const UDPHeaderLen = 8

func getSrcUDPPort() (int, error) {
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		return 0, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).Port, nil
}

func GetUDPHeaderBytes(srcIP, dstIP net.IP, dstPort uint16) ([]byte, error) {

	srcPort, err := getSrcUDPPort()
	if err != nil {
		return nil, fmt.Errorf("get src udp port: %w", err)
	}

	b := make([]byte, UDPHeaderLen)
	binary.BigEndian.PutUint16(b[0:2], uint16(srcPort)) // source port
	binary.BigEndian.PutUint16(b[2:4], dstPort)         // destination port
	binary.BigEndian.PutUint16(b[4:6], UDPHeaderLen)    // length (header only, no data)
	binary.BigEndian.PutUint16(b[6:8], 0)               // checksum placeholder

	checksum, err := udpChecksum(srcIP, dstIP, b)
	if err != nil {
		return nil, fmt.Errorf("udp checksum: %w", err)
	}
	binary.BigEndian.PutUint16(b[6:8], checksum) // Insert checksum
	return b, nil
}

func udpChecksum(srcIP, dstIP net.IP, data []byte) (uint16, error) {

	src, err := srcIP.To4().MarshalText()
	if err != nil {
		return 0, fmt.Errorf("src IP: %w", err)
	}
	dst, err := dstIP.To4().MarshalText()
	if err != nil {
		return 0, fmt.Errorf("dst IP: %w", err)
	}

	var csum uint32

	// Pseudo-header fields
	csum += (uint32(src[0]) + uint32(src[2])) << 8
	csum += uint32(src[1]) + uint32(src[3])
	csum += (uint32(dst[0]) + uint32(dst[2])) << 8
	csum += uint32(dst[1]) + uint32(dst[3])

	csum += uint32(17)                // Protocol number for UDP
	csum += uint32(len(data))         // UDP length

	// UDP header
	for i := 0; i < len(data)-1; i += 2 {
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[len(data)-1]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum), nil
}
