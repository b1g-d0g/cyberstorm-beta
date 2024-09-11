package core

import (
	"context"
	"sync"
)

// RunGoroutines esegue le goroutine in base al profilo scelto
func RunGoroutines(wg *sync.WaitGroup, ctx context.Context, profile string, attackFunc func(), duration int) {
	numGoroutines := getProfileGoroutines(profile)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Esegui l'attacco finché il contesto non è cancellato
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

// getProfileGoroutines restituisce il numero di goroutine in base al profilo
func getProfileGoroutines(profile string) int {
	switch profile {
	case "light":
		return 15
	case "medium":
		return 50
	case "extreme":
		return 100
	default:
		return 50
	}
}
