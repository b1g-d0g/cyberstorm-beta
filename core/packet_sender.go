package core

import (
	"fmt"
	"net"
)

func SendPacket(rawSocket RawSocket, srcIp, dstIp net.IP, protocol int, data []byte) error {
	ipv4Header := GetIPV4Header(srcIp, dstIp, len(data), protocol)
	ipv4HeaderBytes, _ := ipv4Header.Marshal()

	packet := append(ipv4HeaderBytes, data...)

	if err := rawSocket.Send(packet); err != nil {
		return fmt.Errorf("send data: %w", err)
	}
	return nil
}
