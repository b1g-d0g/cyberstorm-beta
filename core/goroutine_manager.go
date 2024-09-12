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

	// Timer per la durata complessiva dell'attacco
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) { // Aggiungi un ID per tracciare la goroutine
			defer wg.Done()
			log.Printf("Goroutine %d iniziata\n", goroutineID) // Log di inizio goroutine

			for {
				select {
				case <-ctx.Done(): // Se il contesto viene cancellato, la goroutine termina
					log.Printf("Goroutine %d terminata (context)\n", goroutineID)
					return
				case <-timer.C: // Se il timer scade, la goroutine termina
					log.Printf("Goroutine %d terminata (timer)\n", goroutineID)
					return
				default:
					attackFunc() // Esegui l'attacco
					time.Sleep(10 * time.Millisecond) // Pausa breve
				}
			}
		}(i) // Passa l'indice della goroutine
	}

	// Aspetta che tutte le goroutine abbiano finito
	wg.Wait()
	log.Println("Tutte le goroutine terminate")
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
		baseGoroutines = 16
	default:
		baseGoroutines = 4
	}

	// Rapporta il numero di goroutine al numero di CPU logiche
	return baseGoroutines * numCPU
}
