# Laboratory Exercise 1

## Description of the problem

The exercise was to create two classes:

- ClearTextFrame
- EncryptedFrame

Which represent a non-encrypted and encrypted frame of the CCMP protocol which
is used in IEEE 802.11i standard (WPA2).

## Requirements

For the two classes I needed to implement the necessary:

- encryption method, for encrypting the data
- decryption method, for decrypting the data
- MIC calculation
- MIC and authenticity verification

## How it works

1. In the main function we take the users message from standard input and after
   that we generate two 16 byte buffers which are the AES key and the IV.

2. We construct a ClearTextFrame using the message and transform it into a byte
   buffer

3. We encrypt the clearText frame and send in the key and the IV as arguments

- This method first creates an array of the message byte buffer
- Creates and AES object from the CryptoSwift library with the key, the IV for the CTR block mode and padding
- It encrypts the array bytes using the aes.encrypt method
- It creates a data object from the array bytes and returns an EncryptedFrame object with the iv, data and key

4. Then we run the sendFrame method which creates a JSON object of the
   EncryptedFrame which can be shared to the client.
5. The client would then decode the JSON and run the decryptAndVerify method on
   his side with the key as an argument. This decrypts the encrypted data in the
   frame, and verifies the mic. If the MIC is the same there isn't an integrity
   problem with the frame, if it isn't there is an error thrown
