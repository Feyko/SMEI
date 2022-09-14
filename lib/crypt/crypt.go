package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
)

func Encrypt(key, value string) (string, error) {
	b := stringTo32B(key)
	block, err := aes.NewCipher(b)
	if err != nil {
		return "", errors.Wrap(err, "could not create a new cipher")
	}

	cipherText := make([]byte, aes.BlockSize+len(value))

	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", errors.New("could not get random bytes")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(value))

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func stringTo32B(str string) []byte {
	var r []byte
	for len(r) < 32 {
		r = append(r, []byte(str)...)
	}
	return r[:32]
}

func Decrypt(key, value string) (string, error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(value)

	if err != nil {
		return "", errors.Wrap(err, "unable to decode value from base64")
	}

	b := stringTo32B(key)
	block, err := aes.NewCipher(b)
	if err != nil {
		return "", errors.Wrap(err, "could not create a new cipher")
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text was too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}