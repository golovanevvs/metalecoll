package decrypt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const decryptKey = ctxKey("decrypt")

// Decrypt - middleware для декодирования
func Decrypt(privateKeyPath string, lg *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encrypted, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			decryptedBody, err := cryptoChecker(privateKeyPath, lg, encrypted)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), decryptKey, decryptedBody))

			next.ServeHTTP(w, r)
		})
	}
}

func cryptoChecker(privateKeyPath string, lg *logrus.Logger, encryptedBody []byte) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(string(encryptedBody))
	if err != nil {
		return make([]byte, 0), err
	}

	// получение шифрованного ключа
	privateKey, err := getPrivateKey(privateKeyPath)
	if err != nil {
		lg.Fatalf("Ошибка инициализации ключа шифрования: %s", err.Error())
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted)
	if err != nil {
		return make([]byte, 0), err
	}

	return decrypted, nil
}

func getPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	file, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("ошибка декодирования")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
