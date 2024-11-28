import argparse
from PIL import Image

def decode_image(image_path, delimiter='11111111'):
    """
    Decodes a hidden message from an image.
    
    Args:
        image_path (str): Path to the image containing the hidden message.
        delimiter (str): Binary string used as a delimiter to mark the end of the message.

    Returns:
        str: The extracted hidden message.
    """
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
    message_bits = binary_message.split(delimiter)[0]

    # Convert binary message to text
    message = ''.join(chr(int(message_bits[i:i + 8], 2)) for i in range(0, len(message_bits), 8))
    return message


def encode_image(image_path, message, output_path, delimiter='11111111'):
    """
    Encodes a hidden message into an image.
    
    Args:
        image_path (str): Path to the input image.
        message (str): The message to hide in the image.
        output_path (str): Path to save the image with the hidden message.
        delimiter (str): Binary string used to mark the end of the message.

    Returns:
        None
    """
    # Open the image
    image = Image.open(image_path)
    binary_message = ''.join(format(ord(char), '08b') for char in message) + delimiter  # Add a delimiter
    pixels = image.load()

    # Ensure the image is in RGB mode
    if image.mode != 'RGB':
        image = image.convert('RGB')

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
    # Preserve format based on output file extension
    output_format = output_path.split('.')[-1].upper()
    if output_format not in ['PNG', 'BMP', 'JPG', 'JPEG']:
        raise ValueError("Unsupported output format. Use PNG, BMP, or JPG/JPEG.")

    image.save(output_path, format=output_format)
    print("Message encoded and image saved as", output_path)


def main():
    parser = argparse.ArgumentParser(
        description="Tool for encoding and decoding hidden messages in images using LSB steganography."
    )
    subparsers = parser.add_subparsers(dest="command", help="Sub-command to execute (encode or decode)")

    # Subparser for encoding
    encode_parser = subparsers.add_parser("encode", help="Encode a message into an image.")
    encode_parser.add_argument("image_path", type=str, help="Path to the input image (PNG, BMP, JPG/JPEG).")
    encode_parser.add_argument("message", type=str, help="The message to hide in the image.")
    encode_parser.add_argument("output_path", type=str, help="Path to save the image with the hidden message (PNG, BMP, JPG/JPEG).")
    encode_parser.add_argument("--delimiter", type=str, default="11111111", help="Binary delimiter to mark the end of the message.")

    # Subparser for decoding
    decode_parser = subparsers.add_parser("decode", help="Decode a hidden message from an image.")
    decode_parser.add_argument("image_path", type=str, help="Path to the image containing the hidden message (PNG, BMP, JPG/JPEG).")
    decode_parser.add_argument("--delimiter", type=str, default="11111111", help="Binary delimiter used to mark the end of the message.")

    # Parse arguments
    args = parser.parse_args()

    if args.command == "encode":
        encode_image(args.image_path, args.message, args.output_path, args.delimiter)
    elif args.command == "decode":
        hidden_message = decode_image(args.image_path, args.delimiter)
        print("Extracted message:", hidden_message)
    else:
        parser.print_help()


if __name__ == "__main__":
    main()