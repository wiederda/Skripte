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

// Function to encode a message into an image
func encodeImage(imagePath, message, outputPath string) error {
	// Open the image file
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("unable to open image: %v", err)
	}
	defer imgFile.Close()

	// Decode the image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}

	// Convert the message into a binary string
	binaryMessage := stringToBinary(message) + "11111111" // Add delimiter (end marker)

	// Get the image bounds and prepare to modify the image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := img.(*image.RGBA)

	dataIndex := 0

	// Iterate over every pixel to embed the message
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the pixel color
			r, g, b, a := pixels.RGBAAt(x, y).R, pixels.RGBAAt(x, y).G, pixels.RGBAAt(x, y).B, pixels.RGBAAt(x, y).A

			// Embed message bits into the LSB of R, G, B channels
			if dataIndex < len(binaryMessage) {
				// Modify the Red channel
				r = (r & 0xFE) | binaryMessage[dataIndex] - '0'
				dataIndex++

				// Modify the Green channel
				if dataIndex < len(binaryMessage) {
					g = (g & 0xFE) | binaryMessage[dataIndex] - '0'
					dataIndex++
				}

				// Modify the Blue channel
				if dataIndex < len(binaryMessage) {
					b = (b & 0xFE) | binaryMessage[dataIndex] - '0'
					dataIndex++
				}
			}

			// Set the new pixel color
			pixels.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			// Stop if the entire message is encoded
			if dataIndex >= len(binaryMessage) {
				break
			}
		}
		if dataIndex >= len(binaryMessage) {
			break
		}
	}

	// Save the modified image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output image file: %v", err)
	}
	defer outputFile.Close()

	// Encode and save the image
	err = png.Encode(outputFile, img)
	if err != nil {
		return fmt.Errorf("unable to encode image: %v", err)
	}

	fmt.Println("Message encoded and image saved as", outputPath)
	return nil
}

// Helper function to convert a string to binary
func stringToBinary(str string) string {
	var binaryString strings.Builder
	for _, char := range str {
		binaryString.WriteString(fmt.Sprintf("%08b", char))
	}
	return binaryString.String()
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
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return "", fmt.Errorf("unable to decode image: %v", err)
	}

	// Convert the image to RGBA format if it's not already in RGBA format
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		// If the image is not in RGBA format, create a new RGBA image and draw the decoded image on it
		rgbaImg = image.NewRGBA(img.Bounds())
		draw.Draw(rgbaImg, img.Bounds(), img, image.Point{0, 0}, draw.Over)
	}

	// Initialize variables
	binaryMessage := ""
	bounds := rgbaImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Iterate through each pixel and extract the LSBs
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := rgbaImg.RGBAAt(x, y)

			// Extract the LSB of each color channel (R, G, B)
			for i := 0; i < 3; i++ { // R, G, B channels
				binaryMessage += fmt.Sprintf("%d", pixel.R>>i&1)
			}
		}
	}

	// Split binary message by the delimiter (the last 8 bits)
	delimiter := "11111111"
	messageBits := strings.Split(binaryMessage, delimiter)[0]

	// Convert binary message to text
	var message string
	for i := 0; i < len(messageBits); i += 8 {
		// Take 8 bits at a time
		byteStr := messageBits[i : i+8]
		// Convert to ASCII
		val, _ := binaryToDecimal(byteStr)
		message += string(rune(val))
	}

	return message, nil
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

// Main function to accept command-line arguments and call encoding/decoding
func main() {
	// Ensure correct number of arguments for encoding or decoding
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <encode|decode> [arguments...]")
		return
	}

	action := os.Args[1]

	if action == "encode" {
		// Encoding
		if len(os.Args) != 5 {
			fmt.Println("Usage for encoding: go run main.go encode <input_image.png> <message> <output_image.png>")
			return
		}
		imagePath := os.Args[2]
		message := os.Args[3]
		outputPath := os.Args[4]

		err := encodeImage(imagePath, message, outputPath)
		if err != nil {
			fmt.Println("Error encoding image:", err)
		}
	} else if action == "decode" {
		// Decoding
		if len(os.Args) != 3 {
			fmt.Println("Usage for decoding: go run main.go decode <input_image.png>")
			return
		}
		imagePath := os.Args[2]

		// Decode the hidden message
		message, err := decodeImage(imagePath)
		if err != nil {
			fmt.Println("Error decoding image:", err)
		} else {
			fmt.Println("Decoded message:", message)
		}
	} else {
		fmt.Println("Unknown action:", action)
	}
}
