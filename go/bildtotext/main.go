package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func printHelp() {
	fmt.Println("Verwendung: go run main.go <Bildpfad>")
	fmt.Println("Extrahiert Text aus einem Bild mithilfe von Tesseract OCR.")
	fmt.Println()
	fmt.Println("Argumente:")
	fmt.Println("  <Bildpfad>   Der Pfad zu dem Bild, aus dem der Text extrahiert werden soll.")
	fmt.Println()
	fmt.Println("Optionen:")
	fmt.Println("  --help       Zeigt diese Hilfe an.")
}

func main() {
	// Überprüfen, ob keine Argumente oder --help übergeben wurde
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		printHelp()
		return
	}

	// Bildpfad aus den Kommandozeilenparametern extrahieren
	imagePath := os.Args[1]

	// Überprüfen, ob die Datei existiert
	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Fatalf("Die angegebene Bilddatei existiert nicht: %s", imagePath)
	}

	// Generiere einen temporären Dateinamen für die Ausgabedatei
	outputFile := "output.txt"

	// Tesseract ausführen und den erkannten Text in eine Datei schreiben lassen
	cmd := exec.Command("tesseract", imagePath, outputFile)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Fehler beim Ausführen von Tesseract: %v", err)
	}

	// Die generierte Textdatei lesen
	outputText, err := os.ReadFile(outputFile + ".txt")
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Ausgabe: %v", err)
	}

	fmt.Println("Erkannter Text:", string(outputText))

	// Optional: Die temporäre Textdatei nach dem Lesen löschen
	err = os.Remove(outputFile + ".txt")
	if err != nil {
		log.Printf("Fehler beim Löschen der temporären Datei: %v", err)
	} else {
		fmt.Println("Temporäre Textdatei gelöscht.")
	}
}
