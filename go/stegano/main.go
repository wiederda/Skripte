package main

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"github.com/auyer/steganography"
)

func main() {
	// Überprüfen der Argumente
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <input.png> <output.png> <message>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	message := os.Args[3]

	// Eingabedatei öffnen
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer file.Close()

	// PNG-Bild decodieren
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Nachricht einbetten
	buffer := new(bytes.Buffer)
	err = steganography.Encode(buffer, img, []byte(message))
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return
	}

	// Ausgabedatei erstellen
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// Modifiziertes Bild speichern
	_, err = buffer.WriteTo(outFile)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Printf("Message encoded and saved to %s\n", outputFile)

	// Zum Testen: Nachricht wieder dekodieren
	outFile.Seek(0, 0) // Datei zurücksetzen
	img, _, _ = image.Decode(outFile)
	decodedMessage := steganography.Decode(steganography.GetMessageSize(img), img)
	fmt.Println("Decoded message:", string(decodedMessage))
}
