package main

import (
	"fmt"
	"log"
	"os"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// Überprüfen, ob ein Bildpfad als Argument übergeben wurde
	if len(os.Args) < 2 {
		fmt.Println("Verwendung: ocr <Pfad-zum-Bild>")
		fmt.Println("Beispiel: ./ocr-linux bild.png")
		return
	}

	// Bildpfad aus dem ersten Argument lesen
	imagePath := os.Args[1]

	// Tesseract-Client erstellen (Fehlerquelle ist vermutlich hier)
	client := gosseract.NewClient()
	defer client.Close()

	// Bildpfad setzen
	client.SetImage(imagePath)

	// Optional: Sprache setzen (z.B. Deutsch und Englisch)
	// client.SetLanguage("eng", "deu")

	// Text aus dem Bild extrahieren
	text, err := client.Text()
	if err != nil {
		log.Fatalf("Fehler bei der Texterkennung: %v\n", err)
	}

	// Erkannten Text ausgeben
	fmt.Println("Erkannter Text:")
	fmt.Println(text)
}
