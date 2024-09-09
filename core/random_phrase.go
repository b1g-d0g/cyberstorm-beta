package core

import (
	"fmt"
	"math/rand"
)

// Funzione per stampare una frase casuale
func PrintRandomPhrase() {
	phrases := []string{
		"Unleash the storm. The grid will never be the same.",
		"Disrupt the flow. Rewrite the code.",
		"When the storm hits, only chaos prevails.",
		"Silent shadows. Violent storms.",
		"The system trembles. The storm rises.",
		"Break the chain. Burn the signal.",
		"We are the glitch in the machine.",
		"Code whispers. Thunder roars.",
		"Chaos is currency in the digital age.",
		"In the eye of the storm, we control the network.",
	}

	// Seleziona una frase casuale
	randomIndex := rand.Intn(len(phrases))
	fmt.Println(phrases[randomIndex])
}
