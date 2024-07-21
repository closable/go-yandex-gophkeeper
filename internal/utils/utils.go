// Package utils are list usfull procs
package utils

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/closable/go-yandex-gophkeeper/internal/store"
	"github.com/golang-jwt/jwt/v4"
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

// StoreFileData store data to file
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

// MakePathBinFile make path
func MakePathFile(path, extension string) string {
	p := strings.Split(path, string(os.PathSeparator))
	if len(p) > 1 {
		p[len(p)-1] = strings.Split(p[len(p)-1], ".")[0] + fmt.Sprintf(".%s", extension)
		return strings.Join(p, string(os.PathSeparator))
	}
	return ""
}

// ZipFolder zip folder information
func ZipFolder(path string) (string, error) {
	output := MakePathFile(path, "zip")
	f := strings.Split(path, string(os.PathSeparator))
	nameFolder := f[len(f)-1]

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(p string, info os.FileInfo, err error) error {
		// fmt.Printf("Crawling: %#v\n", p)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file.Close()
		// f, err := w.Create(p)
		zipPath := p[strings.Index(p, nameFolder):]
		f, err := w.Create(zipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	// find files each folder
	err = filepath.Walk(path, walker)
	if err != nil {
		panic(err)
	}

	return output, nil
}

func GetRandomString(n int) string {
	s := "qwertyuiopasdfghjklzxcvbnm"
	result := ""
	for i := 0; i < n; i++ {
		result += string(s[mrand.Intn(len(s))])
	}
	return result
}

const TokenEXP = time.Hour * 3
const SecretKEY = "*Hello-World*"

// Структура описания JWT токена
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// Функция строитель строки JWT токена, с использованием ID пользователя
func BuildJWTString(userID int) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// Функция парсер JWT токена, для получения ID пользователя
func GetUserID(tokenString string) int {
	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	// time.Sleep(time.Second * 4)
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKEY), nil
	})

	// возвращаем ID пользователя в читаемом виде
	return claims.UserID
}

// OutputListCli output table data
func OutputListCli(data []store.RowItem, decrypt bool, keyString string) {
	fmt.Println(strings.Repeat("_", 153))
	fmt.Printf("| %3s | %-15s | %-55s | %-55s| %-10s |\n", "ИД", "Тип", "Метка", "Данные", "Размер,байт")
	fmt.Println(strings.Repeat("-", 153))

	for _, v := range data {
		data := v.EncData
		name := catStringData(v.Name, 50)
		if decrypt && v.DataType < 3 && len(keyString) > 0 {
			data = Decrypt(keyString, v.EncData)
		}
		data = catStringData(data, 50)
		fmt.Printf("| %3d | %-15s | %-55s | %-55s| %11d |\n", v.Id, v.Type, name, data, v.Length)
	}
	fmt.Println(strings.Repeat("-", 153))
}

// catStringData helper for formatting string
func catStringData(s string, n int) string {
	if len([]byte(s)) > n {
		return s[:n+1] + "..."
	}
	return s
}
