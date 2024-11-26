from PIL import Image


def encode_image(image_path, message, output_path):
    # Open the image
    image = Image.open(image_path)
    binary_message = ''.join(format(ord(char), '08b') for char in message) + '11111111'  # Add a delimiter
    pixels = image.load()

    # Iterate through each pixel and hide the message in the LSB
    data_index = 0
    for y in range(image.height):
        for x in range(image.width):
            pixel = list(pixels[x, y])  # Get the current pixel (R, G, B)

            # Modify the red, green, and blue channels with the message bits
            for i in range(3):  # For each color channel (R, G, B)
                if data_index < len(binary_message):
                    pixel[i] = pixel[i] & 0xFE | int(binary_message[data_index])  # Set LSB to message bit
                    data_index += 1

            # Place the modified pixel back
            pixels[x, y] = tuple(pixel)

            if data_index >= len(binary_message):
                break

    # Save the modified image with hidden message
    image.save(output_path)
    print("Message encoded and image saved as", output_path)


# Example usage
encode_image("input_image.png", 'Dies ist ein Test', "output_image.png")
