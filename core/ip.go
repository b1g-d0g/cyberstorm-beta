package core

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
)

func GetRandPublicIP() net.IP {
	for {
		ip := net.IP(make([]byte, 4))
		binary.BigEndian.PutUint32(ip[0:], uint32(rand.Intn(1<<32-1)))

		// Escludiamo gli IP privati
		if ip[0] != 10 && (ip[0] != 172 || ip[1] < 16 || ip[1] > 31) && ip[0] != 192 && ip[1] != 168 && ip[0] != 169 && ip[1] != 254 {
			return ip
		}
	}
}

func GetIPV4Header(srcIp, dstIp net.IP, dataLen, protocol int) *ipv4.Header {
	return &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TotalLen: ipv4.HeaderLen + dataLen,
		ID:       rand.Intn(65535), // Usa un ID casuale per ogni pacchetto
		TTL:      255,
		Protocol: protocol,
		Src:      srcIp.To4(),
		Dst:      dstIp.To4(),
	}
}
