package core

import (
	"net"
	"syscall"
	"unsafe"
)

type RawSocket struct {
	fd   int
	addr syscall.SockaddrInet4
}

// Crea un raw socket per inviare pacchetti con l'header IP incluso
func NewRawSocket(dstIP net.IP) (RawSocket, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return RawSocket{}, err
	}

	var addr syscall.SockaddrInet4
	copy(addr.Addr[:], dstIP.To4())

	rs := RawSocket{fd: fd, addr: addr}

	// Include l'header IP nei pacchetti inviati
	if err := syscall.SetsockoptInt(rs.fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1); err != nil {
		rs.Close() // Chiude il socket in caso di errore
		return RawSocket{}, err
	}

	return rs, nil
}

// Invia pacchetti utilizzando il raw socket
func (r RawSocket) Send(data []byte) error {
	return syscall.Sendto(r.fd, data, 0, (*syscall.SockaddrInet4)(unsafe.Pointer(&r.addr)))
}

// Chiusura del raw socket
func (r RawSocket) Close() error {
	return syscall.Close(r.fd)
}
