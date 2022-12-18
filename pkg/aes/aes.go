package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	mathrand "math/rand"
	"time"
)

func GenerateKey(keySize int) string {
	// 16 bytes - AES-128
	// 24 bytes - AES-192
	// 32 bytes - AES-256

	mathrand.Seed(time.Now().UnixNano())

	bytes := make([]byte, keySize)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	key := hex.EncodeToString(bytes)
	return key
}

func Encrypt(data []byte, aesKey string) (encryptedString string) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		panic(err)
	}

	plainText := data

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainText, nil)

	return fmt.Sprintf("%x", cipherText)
}

func Decrypt(data []byte, aesKey string) (decryptedString string) {
	key, err := hex.DecodeString(aesKey)
	if err != nil {
		panic(err)
	}
	enc, err := hex.DecodeString(string(data))
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := aesGCM.NonceSize()

	nonce, encryptedText := enc[:nonceSize], enc[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		panic(err)
	}

	return string(plainText)
}
