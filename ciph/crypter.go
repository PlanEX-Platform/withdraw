package ciph

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/zhooq/go-ethereum/crypto/sha3"
)

func Encrypt(plainText string, key string) (cipherText string, nonce string, err error) {
	// подготовка входных параметров
	keyDigest := []byte(key)
	keyNorm := sha3.Sum256(keyDigest)
	dataDigest := []byte(plainText)

	// шифрование
	block, err := aes.NewCipher(keyNorm[:])
	if err != nil {
		return "", "", err
	}
	nonceDigest := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonceDigest)
	if err != nil {
		return "", "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	cipher := aesgcm.Seal(nil, nonceDigest, dataDigest, nil)

	// подготовка выходных параметров
	ciphertext := hex.EncodeToString(cipher)
	nonce = hex.EncodeToString(nonceDigest)
	return ciphertext, nonce, err
}

func Decrypt(cipherText string, nonce string, key string) (plainText string, err error) {
	// подготовка входных параметров
	keyDigest := []byte(key)
	keyNorm := sha3.Sum256(keyDigest)
	cipherDigest, err := hex.DecodeString(cipherText)
	nonceDigest, err := hex.DecodeString(nonce)

	// расшифровка
	block, err := aes.NewCipher(keyNorm[:])
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	textDigest, err := aesgcm.Open(nil, nonceDigest, cipherDigest, nil)
	if err != nil {
		return "", err
	}

	// подготовка выходных параметров
	plainText = string(textDigest)
	return plainText, err
}
