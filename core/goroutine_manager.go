package core

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// RunGoroutines gestisce l'esecuzione delle goroutine in base al profilo.
// La funzione `attackFunc` viene passata come parametro per specificare la logica dell'attacco.
func RunGoroutines(wg *sync.WaitGroup, ctx context.Context, profile string, attackFunc func(), duration int) {
	var numGoroutines int
	var delayBetweenRequests time.Duration

	// Configura il numero di goroutine e il delay tra le richieste in base al profilo
	switch profile {
	case "light":
		numGoroutines = 10
		delayBetweenRequests = 500 * time.Millisecond
		runtime.GOMAXPROCS(1) // Usa una CPU
	case "medium":
		numGoroutines = 50
		delayBetweenRequests = 100 * time.Millisecond
		runtime.GOMAXPROCS(2) // Usa due CPU
	case "extreme":
		numGoroutines = 100
		delayBetweenRequests = 10 * time.Millisecond
		runtime.GOMAXPROCS(runtime.NumCPU()) // Usa tutte le CPU disponibili
	default:
		log.Fatalf("Unknown profile: %s", profile)
	}

	// Avvia le goroutine in base al profilo
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go runSingleGoroutine(wg, ctx, attackFunc, delayBetweenRequests)
		log.Printf("Started goroutine %d with profile %s", i, profile)
	}
}

// runSingleGoroutine esegue una singola goroutine che invoca la funzione dell'attacco
func runSingleGoroutine(wg *sync.WaitGroup, ctx context.Context, attackFunc func(), delayBetweenRequests time.Duration) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("Canceling attack")
			return
		default:
			// Esegue la funzione di attacco
			attackFunc()

			// Attendi il delay tra una richiesta e l'altra
			time.Sleep(delayBetweenRequests)
		}
	}
}
