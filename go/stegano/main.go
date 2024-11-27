package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

// Function to display help
func displayHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go <action> [arguments...]")
	fmt.Println()
	fmt.Println("Actions:")
	fmt.Println("  encode <input_image.png> <message> <output_image.png>")
	fmt.Println("      Encodes the given message into the specified image and saves it.")
	fmt.Println()
	fmt.Println("  decode <input_image.png>")
	fmt.Println("      Decodes and prints the hidden message from the specified image.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help")
	fmt.Println("      Displays this help message.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  Encode a message:")
	fmt.Println("      go run main.go encode input.png \"Hidden message\" output.png")
	fmt.Println()
	fmt.Println("  Decode a message:")
	fmt.Println("      go run main.go decode output.png")
	fmt.Println()
}

// Function to encode a message into an image
func encodeImage(imagePath, message, outputPath string) error {
	// Open the image file
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("unable to open image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	srcImg, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}

	// Ensure the image is in RGBA format
	bounds := srcImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, srcImg, image.Point{}, draw.Src)

	// Convert the message into a binary string
	binaryMessage := stringToBinary(message) + "11111111" // Add delimiter (end marker)

	// Embed the binary message
	dataIndex := 0
	for y := bounds.Min.Y; y < bounds.Max.Y && dataIndex < len(binaryMessage); y++ {
		for x := bounds.Min.X; x < bounds.Max.X && dataIndex < len(binaryMessage); x++ {
			// Get current pixel
			r, g, b, a := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Embed bits into the LSBs of R, G, B channels
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

			// Set modified pixel
			img.Set(x, y, color.RGBA{r8, g8, b8, uint8(a >> 8)})
		}
	}

	// Save the modified image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output image file: %v", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		return fmt.Errorf("unable to encode image: %v", err)
	}

	fmt.Println("Message encoded and image saved as", outputPath)
	return nil
}

// Function to decode a message from an image
func decodeImage(imagePath string) (string, error) {
	// Open the image file
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("unable to open image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	srcImg, _, err := image.Decode(imgFile)
	if err != nil {
		return "", fmt.Errorf("unable to decode image: %v", err)
	}

	// Ensure the image is in RGBA format
	bounds := srcImg.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, srcImg, image.Point{}, draw.Src)

	// Extract binary message
	var binaryMessage strings.Builder
	delimiter := "11111111"
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Append LSBs of R, G, B channels
			binaryMessage.WriteByte(r8&1 + '0')
			binaryMessage.WriteByte(g8&1 + '0')
			binaryMessage.WriteByte(b8&1 + '0')

			// Check if delimiter is reached
			if strings.HasSuffix(binaryMessage.String(), delimiter) {
				goto EndDecoding
			}
		}
	}
EndDecoding:

	// Remove the delimiter and convert to text
	messageBits := binaryMessage.String()
	messageBits = strings.TrimSuffix(messageBits, delimiter)

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

// Helper function to convert a string to binary
func stringToBinary(str string) string {
	var binaryString strings.Builder
	for _, char := range str {
		binaryString.WriteString(fmt.Sprintf("%08b", char))
	}
	return binaryString.String()
}

// Helper function to convert binary string to decimal
func binaryToDecimal(binStr string) (int, error) {
	var decimalValue int
	for i, bit := range binStr {
		if bit == '1' {
			decimalValue += 1 << (7 - i)
		}
	}
	return decimalValue, nil
}

// Main function
func main() {
	// Check if the user needs help
	if len(os.Args) < 2 || os.Args[1] == "--help" {
		displayHelp()
		return
	}

	// Parse the action (encode/decode)
	action := os.Args[1]

	if action == "encode" {
		// Encoding
		if len(os.Args) != 5 {
			fmt.Println("Error: Incorrect number of arguments for 'encode'.")
			fmt.Println("Usage: go run main.go encode <input_image.png> <message> <output_image.png>")
			return
		}
		err := encodeImage(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			fmt.Println("Error encoding image:", err)
		}
	} else if action == "decode" {
		// Decoding
		if len(os.Args) != 3 {
			fmt.Println("Error: Incorrect number of arguments for 'decode'.")
			fmt.Println("Usage: go run main.go decode <input_image.png>")
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
