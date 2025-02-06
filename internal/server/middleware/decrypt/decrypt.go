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

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/sirupsen/logrus"
)

// Decrypt - middleware для декодирования
func Decrypt(privateKeyPath string, lg *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encrypted, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			lg.Errorf("Ошибка чтения зашифрованного тела: %s", err.Error())
			return
		}

		decryptedBody, err := cryptoChecker(privateKeyPath, lg, encrypted)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			lg.Errorf("Ошибка расшифровки зашифрованного тела: %s", err.Error())
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), constants.DecryptKey, decryptedBody))

		fmt.Println("Из мидлварь", decryptedBody)

		next.ServeHTTP(w, r)
	})
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
		return nil, fmt.Errorf("не удалось прочитать файл: %s", err.Error())
	}

	block, rest := pem.Decode(file)
	if len(rest) > 0 {
		return nil, fmt.Errorf("ошибка декодирования: в файл обнаружены непредвиденные данные после PEM-блока")
	}
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		fmt.Printf("block: %v, block.Type: %v", block, block.Type)
		return nil, fmt.Errorf("ошибка декодирования: неверный формат файла")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования: не удалось распарсить приватный ключ: %s", err.Error())
	}

	return privateKey, nil
}
