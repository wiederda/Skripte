package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Kommandozeilenparameter definieren
	url := flag.String("url", "https://example.com oder https://www.example.com", "URL dessen Zertifikat abgefragt werden soll")
	outputFile := flag.String("output", "cert_info.txt", "Output File in das die Zertifikatsinformationen geschrieben werden")

	// Parse the Kommandozeilenparameter
	flag.Parse()

	// Überprüfen, ob die URL angegeben wurde
	if *url == "" {
		fmt.Println("Please provide a URL using the --url flag")
		flag.PrintDefaults()
		return
	}

	// Datei zum Schreiben öffnen
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Writer für UTF-8-Encoding erstellen
	writer := bufio.NewWriter(file)

	// HTTP Transport einstellen
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}

	// Anfrage stellen
	resp, err := client.Get(*url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Zertifikat ermitteln
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]

		// Zertifikatsinformationen in die Datei schreiben
		fmt.Fprintf(writer, "Zertifikat für: %s\n", cert.Subject.CommonName)
		fmt.Fprintf(writer, "Ausgestellt von: %s\n", cert.Issuer.CommonName)
		fmt.Fprintf(writer, "Gültig von: %s\n", cert.NotBefore)
		fmt.Fprintf(writer, "Ablaufdatum: %s\n", cert.NotAfter)
		fmt.Println("Zertifikatsinformationen wurden in die Datei geschrieben:", *outputFile)
	} else {
		fmt.Fprintln(writer, "Keine TLS-Verbindung oder kein Zertifikat gefunden")
		fmt.Println("Keine TLS-Verbindung oder kein Zertifikat gefunden")
	}

	// Writer flushen, um sicherzustellen, dass alle Daten geschrieben werden
	writer.Flush()
}
