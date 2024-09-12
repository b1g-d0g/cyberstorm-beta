package core

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// RunGoroutines esegue le goroutine in base al profilo scelto e in proporzione al numero di CPU logiche disponibili
func RunGoroutines(wg *sync.WaitGroup, ctx context.Context, profile string, attackFunc func(), duration int) {
	numGoroutines := getProfileGoroutines(profile)

	// Utilizza un contesto con timeout per gestire la durata complessiva
	ctx, cancel := context.WithTimeout(ctx, time.Duration(duration)*time.Second)
	defer cancel()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {  // Aggiungi un ID per tracciare la goroutine
			defer wg.Done()
			log.Printf("Goroutine %d iniziata\n", goroutineID) // Log di inizio goroutine

			for {
				select {
				case <-ctx.Done(): // Interrompi l'esecuzione quando il contesto scade
					log.Printf("Goroutine %d terminata\n", goroutineID) // Log di fine goroutine
					return
				default:
					attackFunc() // Esegui l'attacco
					time.Sleep(10 * time.Millisecond) // Pausa breve
				}
			}
		}(i) // Passa l'indice della goroutine
	}
}

// getProfileGoroutines restituisce il numero di goroutine in base al profilo e alle CPU logiche disponibili
func getProfileGoroutines(profile string) int {
	// Determina il numero di CPU logiche disponibili
	numCPU := runtime.NumCPU()

	// Determina il numero base di goroutine per ciascun profilo
	var baseGoroutines int
	switch profile {
	case "light":
		baseGoroutines = 4
	case "medium":
		baseGoroutines = 8
	case "extreme":
		baseGoroutines = 12
	default:
		baseGoroutines = 4
	}

	// Rapporta il numero di goroutine al numero di CPU logiche
	return baseGoroutines * numCPU
}

