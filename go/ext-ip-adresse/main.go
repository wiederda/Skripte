package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// URL eines externen Dienstes, der die öffentliche IP-Adresse zurückgibt
	url := "https://api.ipify.org?format=json"

	// HTTP-Anfrage an den Dienst senden
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Fehler beim Abrufen der IP-Adresse:", err)
		return
	}
	defer resp.Body.Close()

	// Antwort auslesen
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fehler beim Lesen der Antwort:", err)
		return
	}

	// JSON-Antwort parsen
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Fehler beim Parsen des JSON:", err)
		return
	}

	// Öffentliche IP-Adresse ausgeben
	fmt.Println("Zugewiesene öffentliche IP-Adresse des Internetanbieters:", result["ip"])
}
