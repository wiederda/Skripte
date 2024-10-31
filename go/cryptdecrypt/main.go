package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const (
	saltLength = 16 // Länge des Salts
	keyLength  = 32 // AES-256
)

// Schlüsselableitung mit scrypt
func deriveKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 1<<15, 8, 1, keyLength)
}

// Padding für AES
func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Entpadded die Daten
func unpad(data []byte) ([]byte, error) {
	padding := data[len(data)-1]
	if int(padding) > aes.BlockSize {
		return nil, fmt.Errorf("ungültiges Padding")
	}
	return data[:len(data)-int(padding)], nil
}

// Funktion zur AES-Verschlüsselung
func encrypt(plaintext []byte, password string) (string, string, error) {
	// Erzeuge ein zufälliges Salt
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", "", err
	}

	// Leite den Schlüssel ab
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", "", err
	}

	// Erstelle einen neuen AES-Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// Padding hinzufügen
	plaintext = pad(plaintext)

	// Erstelle einen IV (Initialisierungsvektor)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	// Verschlüssele
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Rückgabe von Salt und Ciphertext als Base64-kodierte Strings
	return base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Funktion zur AES-Entschlüsselung
func decrypt(saltBase64, ciphertextBase64, password string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, err
	}

	// Leite den Schlüssel ab
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	// Erstelle einen neuen AES-Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Entschlüssele
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("Ciphertext zu kurz")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Entpaddede Daten zurückgeben
	return unpad(ciphertext)
}

func main() {
	// Kommandozeilenflags für den Betriebsmodus, Text und Passwort
	mode := flag.String("mode", "", "Betriebsmodus: 'crypt' für Verschlüsselung oder 'decrypt' für Entschlüsselung")
	text := flag.String("text", "", "Text zum Verschlüsseln/Entschlüsseln (im Format 'salt:ciphertext' für Decrypt)")
	password := flag.String("password", "", "Passwort für Verschlüsselung/Entschlüsselung")

	flag.Parse()

	// Überprüfe, ob erforderliche Parameter angegeben sind
	if *mode == "" || *password == "" {
		fmt.Println("Fehler: Alle Parameter (mode, password) müssen angegeben werden.")
		flag.Usage()
		os.Exit(1)
	}

	if *mode == "crypt" {
		// Überprüfe, ob der Text für die Verschlüsselung angegeben wurde
		if *text == "" {
			fmt.Println("Fehler: Text zum Verschlüsseln muss angegeben werden.")
			flag.Usage()
			os.Exit(1)
		}

		// Verschlüsselung
		salt, ciphertext, err := encrypt([]byte(*text), *password)
		if err != nil {
			log.Fatalf("Fehler bei der Verschlüsselung: %v", err)
		}

		fmt.Printf("Salt: %s\n", salt)
		fmt.Printf("Ciphertext: %s\n", ciphertext)
	} else if *mode == "decrypt" {
		// Überprüfe, ob der Text für die Entschlüsselung angegeben wurde
		if *text == "" {
			fmt.Println("Fehler: Text im Format 'salt:ciphertext' muss angegeben werden.")
			flag.Usage()
			os.Exit(1)
		}

		// Entschlüsselung erwartet Salt und Ciphertext im Textfeld
		parts := splitSaltCiphertext(*text)
		if len(parts) != 2 {
			log.Fatalf("Text muss im Format 'salt:ciphertext' für die Entschlüsselung sein.")
		}

		salt, ciphertext := parts[0], parts[1]
		decryptedText, err := decrypt(salt, ciphertext, *password)
		if err != nil {
			log.Fatalf("Fehler bei der Entschlüsselung: %v", err)
		}

		// Nur den entschlüsselten Text ausgeben
		fmt.Print(string(decryptedText))
	} else {
		fmt.Println("Ungültiger Modus. Verwenden Sie 'crypt' oder 'decrypt'.")
		flag.Usage()
		os.Exit(1)
	}
}

// Hilfsfunktion zum Aufteilen von Salt und Ciphertext
func splitSaltCiphertext(input string) []string {
	parts := strings.SplitN(input, ":", 2) // Teilen in zwei Teile
	if len(parts) != 2 {
		log.Fatalf("Text muss im Format 'salt:ciphertext' sein.")
	}
	return parts
}