package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/zenazn/pkcs7pad"
	"io"
)

type EncryptedMessage struct {
	Cyphertext  []byte `json:"cyphertext"`
	MessageSize int    `json:"message_size"`
	IV          []byte `json:"iv"`
}

type PlainMessage struct {
	Message string `json:"message"`
}

func (plaintext *PlainMessage) Encrypt() (EncryptedMessage, error) {
	_, err := generateKey()
	if err != nil {
		panic(err)
	}

	message := []byte(plaintext.Message)
	message = pkcs7pad.Pad(message, aes.BlockSize)

	block, err := aes.NewCipher([]byte("random"))
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(message))

	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], message)

	return EncryptedMessage{
		Cyphertext:  ciphertext,
		MessageSize: len(message),
		IV:          iv,
	}, nil
}

func (encrypted *EncryptedMessage) Decrypt() PlainMessage {
	plaintext := make([]byte, encrypted.MessageSize)
	block, err := aes.NewCipher([]byte("random"))
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, encrypted.IV)
	stream.XORKeyStream(plaintext, encrypted.Cyphertext[aes.BlockSize:])

	return PlainMessage{string(plaintext)}
}
