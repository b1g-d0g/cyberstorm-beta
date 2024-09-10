package core

import (
	"context"
	"log"
	"sync"
)

// RunGoroutines esegue un numero di goroutine in base al profilo scelto
func RunGoroutines(wg *sync.WaitGroup, ctx context.Context, profile string, attackFunc func(), duration int) {
	numGoroutines := getNumGoroutines(profile)

	log.Printf("Starting attack with %d goroutines for %d seconds\n", numGoroutines, duration)

	for i := 0; numGoroutines > i; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					attackFunc()
				}
			}
		}()
	}
}

// getNumGoroutines restituisce il numero di goroutine in base al profilo selezionato
func getNumGoroutines(profile string) int {
	switch profile {
	case "light":
		return 10
	case "medium":
		return 50
	case "extreme":
		return 200
	default:
		return 50
	}
}
