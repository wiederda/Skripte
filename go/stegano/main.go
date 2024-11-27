package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	_ "golang.org/x/image/bmp" // Import für BMP-Unterstützung
)

// Funktion zur Anzeige der Hilfe
func displayHelp() {
	fmt.Println("Usage:")
	fmt.Println("  ./stegano <action> [arguments...]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  encode <input_image> <message> <output_image>")
	fmt.Println("      Encodes the given message into the specified image and saves it.")
	fmt.Println()
	fmt.Println("  decode <input_image>")
	fmt.Println("      Decodes and prints the hidden message from the specified image.")
	fmt.Println()
	fmt.Println("Supported Formats:")
	fmt.Println("  Input: PNG, BMP, JPG/JPEG")
	fmt.Println("  Output: PNG, BMP, JPG/JPEG")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help")
	fmt.Println("      Displays this help message.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  Encode a message:")
	fmt.Println("      ./stegano encode input.jpg \"Hidden message\" output.png")
	fmt.Println()
	fmt.Println("  Decode a message:")
	fmt.Println("     ./stegano decode output.png")
	fmt.Println()
}

// Funktion zum Speichern des Bildes im richtigen Format
func saveImage(outputPath string, img image.Image) error {
	ext := strings.ToLower(filepath.Ext(outputPath))
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output image file: %v", err)
	}
	defer outputFile.Close()

	switch ext {
	case ".png":
		err = png.Encode(outputFile, img)
	case ".jpg", ".jpeg":
		// JPEG mit Qualität 90 speichern
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 90})
	case ".bmp":
		// BMP wird durch die go-x-image/bmp unterstützt
		err = bmp.Encode(outputFile, img)
	default:
		return fmt.Errorf("unsupported output format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("unable to encode image: %v", err)
	}
	return nil
}

// Funktion zum Einbetten einer Nachricht in ein Bild
func encodeImage(imagePath, message, outputPath string) error {
	// Bilddatei öffnen
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("unable to open image: %v", err)
	}
	defer imgFile.Close()

	// Bild dekodieren
	srcImg, format, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}
	fmt.Printf("Input image format: %s\n", format)

	// In RGBA-Format umwandeln
	bounds := srcImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, srcImg, image.Point{}, draw.Src)

	// Nachricht in Binärdaten umwandeln
	binaryMessage := stringToBinary(message) + "11111111" // Delimiter hinzufügen

	// Binärnachricht einbetten
	dataIndex := 0
	for y := bounds.Min.Y; y < bounds.Max.Y && dataIndex < len(binaryMessage); y++ {
		for x := bounds.Min.X; x < bounds.Max.X && dataIndex < len(binaryMessage); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// LSBs ändern
			if dataIndex < len(binaryMessage) {
				r8 = (r8 & 0xFE) | (binaryMessage[dataIndex] - '0')
				dataIndex++
			}
			if dataIndex < len(binaryMessage) {
				g8 = (g8 & 0xFE) | (binaryMessage[dataIndex] - '0')
				dataIndex++
			}
			if dataIndex < len(binaryMessage) {
				b8 = (b8 & 0xFE) | (binaryMessage[dataIndex] - '0')
				dataIndex++
			}

			// Pixel setzen
			img.Set(x, y, color.RGBA{r8, g8, b8, uint8(a >> 8)})
		}
	}

	// Bild speichern
	err = saveImage(outputPath, img)
	if err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}

	fmt.Println("Message encoded and image saved as", outputPath)
	return nil
}

// Funktion zum Dekodieren einer Nachricht aus einem Bild
func decodeImage(imagePath string) (string, error) {
	// Bilddatei öffnen
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("unable to open image: %v", err)
	}
	defer imgFile.Close()

	// Bild dekodieren
	srcImg, _, err := image.Decode(imgFile)
	if err != nil {
		return "", fmt.Errorf("unable to decode image: %v", err)
	}

	// In RGBA-Format umwandeln
	bounds := srcImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, srcImg, image.Point{}, draw.Src)

	// Binäre Nachricht extrahieren
	var binaryMessage strings.Builder
	delimiter := "11111111"
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// LSBs extrahieren
			binaryMessage.WriteByte(r8&1 + '0')
			binaryMessage.WriteByte(g8&1 + '0')
			binaryMessage.WriteByte(b8&1 + '0')

			// Prüfen, ob das Delimiter erreicht ist
			if strings.HasSuffix(binaryMessage.String(), delimiter) {
				goto EndDecoding
			}
		}
	}
EndDecoding:

	// Delimiter entfernen und Binärdaten in Text umwandeln
	messageBits := strings.TrimSuffix(binaryMessage.String(), delimiter)
	var message strings.Builder
	for i := 0; i < len(messageBits); i += 8 {
		if i+8 > len(messageBits) {
			break
		}
		byteStr := messageBits[i : i+8]
		val, _ := binaryToDecimal(byteStr)
		message.WriteByte(byte(val))
	}

	return message.String(), nil
}

// Hilfsfunktionen
func stringToBinary(str string) string {
	var binaryString strings.Builder
	for _, char := range str {
		binaryString.WriteString(fmt.Sprintf("%08b", char))
	}
	return binaryString.String()
}

func binaryToDecimal(binStr string) (int, error) {
	var decimalValue int
	for i, bit := range binStr {
		if bit == '1' {
			decimalValue += 1 << (7 - i)
		}
	}
	return decimalValue, nil
}

// Hauptfunktion
func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		displayHelp()
		return
	}

	action := os.Args[1]

	if action == "encode" {
		if len(os.Args) != 5 {
			fmt.Println("Error: Incorrect number of arguments for 'encode'.")
			fmt.Println("Usage: go run main.go encode <input_image> <message> <output_image>")
			return
		}
		err := encodeImage(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println("Error encoding image:", err)
		}
	} else if action == "decode" {
		if len(os.Args) != 3 {
			fmt.Println("Error: Incorrect number of arguments for 'decode'.")
			fmt.Println("Usage: go run main.go decode <input_image>")
			return
		}
		message, err := decodeImage(os.Args[2])
		if err != nil {
			fmt.Println("Error decoding image:", err)
		} else {
			fmt.Println("Decoded message:", message)
		}
	} else {
		fmt.Println("Error: Unknown action:", action)
		fmt.Println("Use '--help' to see available actions.")
	}
}
