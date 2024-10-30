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

3. An EncryptedFrame is created with the key and iv and we send the
   clearTextFrame object for encryption using the class method encrypt.

   - The Encrypt Method creates an array of bytes from the message and creas and
     AES object using the helper library CryptoSwift, we pass it the key and set
     the bloc mode to CTR along with the iv as well as setting the padding.
   - We use the encrypt function from the library to encrypt the clearTextBytes
     we got from the clearTextFrame object which afterward we set the encryptedData property of
     the EncryptedFrame class to the result of the method
   - After that is set we call the calculateMIC function which calculates a MIC
     using hmac, and the key using SHA256 and the encryptedData. This sets the mic
     property of the EncryptedFrame class to the resulting bytes of the process.

4. After the encryption we run the method `sendFrame` which runs a 'simulation'
   of how a sending of such a frame would look like and just returns a
   dictionary that holds the encryptedData and the mic.

5. We then decrypt this data using the decrypt and verify method of of the
   EncryptedFrame class which we send the receivedData, the receivedMIC and the
   IV.
   - First I convert the receivedData to an array of bytes
   - Create the aes object using the key, iv and set the padding and try to
     decrypt the receivedByte array
   - The decryptedData is then turned into an byte buffer which then we
     calculate the mic using the hmac and check if it's the same with the
     receivedMIC we send, if it's not throw an integrity check error, if it is
     return the data
6. Turn the data into a hexadecimal string representation and a normal string
   and print it to the console.
