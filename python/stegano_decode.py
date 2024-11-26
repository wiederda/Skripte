from PIL import Image

def decode_image(image_path):
    image = Image.open(image_path)
    pixels = image.load()
    binary_message = ""

    # Iterate through each pixel and extract the LSBs
    for y in range(image.height):
        for x in range(image.width):
            pixel = pixels[x, y]
            for i in range(3):  # For each color channel (R, G, B)
                binary_message += str(pixel[i] & 1)  # Extract the LSB

    # Split binary message by the delimiter (the last 8 bits)
    delimiter = '11111111'
    message_bits = binary_message.split(delimiter)[0]

    # Convert binary message to text
    message = ''.join(chr(int(message_bits[i:i + 8], 2)) for i in range(0, len(message_bits), 8))
    return message


# Example usage
hidden_message = decode_image("output_image.png")
print("Extracted message:", hidden_message)