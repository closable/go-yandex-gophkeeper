// Package utils are list usfull procs
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Encrypt use for encrypting data
func Encrypt(keyString, stringToEncrypt string) string {
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)

}

// Decrypt use for decrypting data
func Decrypt(keyString, stringToDecrypt string) string {
	key, _ := hex.DecodeString(keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext[:])

}

// CryptoSeq helper function for creating crypto sequence
func CryptoSeq(s string) (string, error) {
	key := []byte(s)

	// add symbol * until 32 bit if len(s) less 32 bit
	if len(s) < 32 {
		key = []byte(s + strings.Repeat("*", 32-len(s)))
	}

	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// GetFileData get data from file
func GetFileData(path string) ([]byte, error) {
	data := make([]byte, 0)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return data, errors.New("file not exists")
	}
	defer file.Chdir()

	data, err = os.ReadFile(path)

	return data, err

}

func StoreFileData(path string, cryptoText string) error {

	file, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	_, err = file.WriteString(cryptoText)
	if err != nil {
		return err
	}
	return nil
}

func MakePathBinFile(path string) string {
	p := strings.Split(path, string(os.PathSeparator))
	if len(p) > 1 {
		p[len(p)-1] = strings.Split(p[len(p)-1], ".")[0] + ".bin"
		return strings.Join(p, string(os.PathSeparator))
	}
	return ""
}
