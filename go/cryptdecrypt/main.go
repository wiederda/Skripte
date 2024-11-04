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
	saltLength = 16 // Length of the salt
	keyLength  = 32 // AES-256 key length
)

// DeriveKey uses scrypt to derive a key from the password and salt
func deriveKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 1<<15, 8, 1, keyLength)
}

// Pad applies PKCS7 padding to data
func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Unpad removes PKCS7 padding from data
func unpad(data []byte) ([]byte, error) {
	padding := data[len(data)-1]
	if int(padding) > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-int(padding)], nil
}

// Encrypt encrypts plaintext using AES-256 with a password-derived key
func encrypt(plaintext []byte, password string) (string, string, error) {
	// Generate a random salt
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", "", err
	}

	// Derive the key
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// Pad the plaintext
	plaintext = pad(plaintext)

	// Create IV and ciphertext buffer
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	// Encrypt using CBC mode
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// Return salt and ciphertext as Base64-encoded strings
	return base64.StdEncoding.EncodeToString(salt), base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts Base64-encoded salt and ciphertext using AES-256 with a password-derived key
func decrypt(saltBase64, ciphertextBase64, password string) ([]byte, error) {
	// Decode the Base64-encoded salt
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding salt: %v", err)
	}

	// Decode the Base64-encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding ciphertext: %v", err)
	}

	// Derive the key
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Check ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract IV from ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt using CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad the decrypted data and return it
	return unpad(ciphertext)
}

func main() {
	// Command-line flags for mode, text, and password
	mode := flag.String("mode", "", "Mode: 'crypt' for encryption or 'decrypt' for decryption")
	text := flag.String("text", "", "Text to encrypt/decrypt (format 'salt:ciphertext' for decryption)")
	password := flag.String("password", "", "Password for encryption/decryption")

	flag.Parse()

	// Check if required parameters are provided
	if *mode == "" || *password == "" {
		fmt.Println("Error: All parameters (mode, password) are required.")
		flag.Usage()
		os.Exit(1)
	}

	if *mode == "crypt" {
		// Ensure text is provided for encryption
		if *text == "" {
			fmt.Println("Error: Text to encrypt is required.")
			flag.Usage()
			os.Exit(1)
		}

		// Perform encryption
		salt, ciphertext, err := encrypt([]byte(*text), *password)
		if err != nil {
			log.Fatalf("Encryption error: %v", err)
		}

		// Print salt and ciphertext only for encryption mode
		fmt.Printf("Salt: %s\n", salt)
		fmt.Printf("Ciphertext: %s\n", ciphertext)
	} else if *mode == "decrypt" {
		// Ensure text is provided in 'salt:ciphertext' format for decryption
		if *text == "" {
			fmt.Println("Error: Text in 'salt:ciphertext' format is required.")
			flag.Usage()
			os.Exit(1)
		}

		// Split salt and ciphertext
		parts := splitSaltCiphertext(*text)
		if len(parts) != 2 {
			log.Fatalf("Text must be in 'salt:ciphertext' format for decryption.")
		}

		// Perform decryption
		salt, ciphertext := parts[0], parts[1]
		decryptedText, err := decrypt(salt, ciphertext, *password)
		if err != nil {
			log.Fatalf("Decryption error: %v", err)
		}

		// Output the decrypted text only (no salt/ciphertext for decrypt mode)
		fmt.Print(string(decryptedText))
	} else {
		fmt.Println("Invalid mode. Use 'crypt' or 'decrypt'.")
		flag.Usage()
		os.Exit(1)
	}
}

// splitSaltCiphertext splits input in the format 'salt:ciphertext'
func splitSaltCiphertext(input string) []string {
	parts := strings.SplitN(input, ":", 2) // Split into two parts
	if len(parts) != 2 {
		log.Fatalf("Text must be in 'salt:ciphertext' format.")
	}
	return parts
}
