package main

import (
	"bufio"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Function to read subject information from a file
func readSubjectFromFile(filename string) (pkix.Name, error) {
	file, err := os.Open(filename)
	if err != nil {
		return pkix.Name{}, fmt.Errorf("Fehler beim Öffnen der Datei: %v", err)
	}
	defer file.Close()

	subj := pkix.Name{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return pkix.Name{}, fmt.Errorf("Ungültiges Format in der Datei")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "CommonName":
			subj.CommonName = value
		case "Organization":
			subj.Organization = []string{value}
		case "Country":
			subj.Country = []string{value}
		}
	}

	if err := scanner.Err(); err != nil {
		return pkix.Name{}, fmt.Errorf("Fehler beim Lesen der Datei: %v", err)
	}

	return subj, nil
}

func main() {
	// Flags für den Algorithmus und die Betreffdatei
	algorithm := flag.String("algorithm", "ecdsa", "Algorithmus für den Schlüssel (ecdsa/rsa)")
	subjectFile := flag.String("subject", "subject.txt", "Datei mit den Betreffinformationen")
	flag.Parse()

	// Lese Betreffinformationen aus der Datei
	subj, err := readSubjectFromFile(*subjectFile)
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Betreffinformationen: %v", err)
	}

	var privateKey crypto.PrivateKey

	switch strings.ToLower(*algorithm) {
	case "ecdsa":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalf("Fehler beim Erstellen des ECDSA-Schlüssels: %v", err)
		}
	case "rsa":
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Fatalf("Fehler beim Erstellen des RSA-Schlüssels: %v", err)
		}
	default:
		log.Fatalf("Ungültiger Algorithmus ausgewählt: %s", *algorithm)
	}

	var sigAlg x509.SignatureAlgorithm
	switch privateKey := privateKey.(type) {
	case *ecdsa.PrivateKey:
		sigAlg = x509.ECDSAWithSHA256
	case *rsa.PrivateKey:
		sigAlg = x509.SHA256WithRSA
	default:
		log.Fatalf("Ungültiger Schlüsseltyp: %T", privateKey)
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: sigAlg,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen der CSR: %v", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	var privPEM []byte
	switch privateKey := privateKey.(type) {
	case *ecdsa.PrivateKey:
		privBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			log.Fatalf("Fehler beim Marshallen des ECDSA-Schlüssels: %v", err)
		}
		privPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privBytes,
		})
	case *rsa.PrivateKey:
		privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privBytes,
		})
	}

	err = os.WriteFile("csr.pem", csrPEM, 0644)
	if err != nil {
		log.Fatalf("Fehler beim Speichern der CSR: %v", err)
	}
	err = os.WriteFile("private_key.pem", privPEM, 0644)
	if err != nil {
		log.Fatalf("Fehler beim Speichern des privaten Schlüssels: %v", err)
	}

	fmt.Println("CSR und privater Schlüssel wurden erfolgreich erstellt und gespeichert.")
}
