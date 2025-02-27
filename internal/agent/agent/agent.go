// Модуль agent предназначен для запуска агента.
package agent

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/mapstorage"
	pb "github.com/golovanevvs/metalecoll/internal/proto"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type agent struct {
	store          mapstorage.Storage
	pollInterval   int
	reportInterval int
}

// Start запускает агента.
func Start(config *config) {
	var putString string
	var body Metrics
	//var metricsJSONGZIP bytes.Buffer

	publicCryptoKey, err := getPublicKey(config.PublicKeyPath)
	if err != nil {
		fmt.Printf("Ошибка получения публичного ключа: %s\n", err.Error())
		os.Exit(1)
	}

	store := mapstorage.NewStorage()

	ag := NewAgent(store, config.PollInterval, config.ReportInterval)

	client := &http.Client{}

	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)

	defer pollIntTime.Stop()
	defer reportIntTime.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit
		cancel()
	}(cancel)

	// выбор сервера REST API или gRPC
	server := "gRPC"

	var gRPCClient *grpc.ClientConn
	var gRPCMetricsClient pb.MetricsClient
	if server == "gRPC" {
		gRPCClient, err = grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("ошибка соединения с gRPC сервером: %s\n", err.Error())
			return
		}
		defer gRPCClient.Close()

		gRPCMetricsClient = pb.NewMetricsClient(gRPCClient)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollIntTime.C:
			RegisterMetrics(ag)
		case <-reportIntTime.C:
			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("Reporting...")

			fmt.Println("Получение данных из хранилища...")
			mapStore, err := ag.store.GetMetricsMap()
			if err != nil {
				fmt.Println("Ошибка получения данных из хранилища:", err)
				continue
			}
			fmt.Println(mapStore)
			fmt.Println("Получение данных из хранилища прошло успешно")

			fmt.Println("Формирование метрики...")

			for _, value := range mapStore {
				switch value.Type {
				case constants.GaugeType:
					v, _ := value.Value.(float64)
					body = Metrics{
						ID:    value.Name,
						MType: value.Type,
						Value: &v,
					}
				case constants.CounterType:
					v, _ := value.Value.(int64)
					body = Metrics{
						ID:    value.Name,
						MType: value.Type,
						Delta: &v,
					}
				}

				metrics := body

				fmt.Println("Формирование метрики прошло успешно")
				fmt.Println(metrics)

				// отправка метрик с помощью REST API или gRPC
				switch server {

				case "REST API":

					putString = fmt.Sprintf("http://%s/update/", config.Addr)

					fmt.Println("Кодирование в JSON...")
					metricsJSON, err := json.Marshal(metrics)
					if err != nil {
						fmt.Println("Ошибка кодирования в JSON:", err)
						continue
					}
					fmt.Println("Кодирование в JSON прошло успешно")

					// fmt.Println("Сжатие в gzip...")
					// gzipWr := gzip.NewWriter(&metricsJSONGZIP)
					// _, err = gzipWr.Write(metricsJSON)
					// if err != nil {
					// 	fmt.Println("Ошибка сжатия в gzip:", err)
					// 	gzipWr.Close()
					// 	continue
					// }
					// gzipWr.Close()
					// fmt.Println("Сжатие в gzip прошло успешно")

					fmt.Println("encryptedMessageBase64")
					encryptedMessageBase64, err := encryptBody(metricsJSON, publicCryptoKey)
					if err != nil {
						fmt.Println("Ошибка кодирования в base64:", err)
						continue
					}

					var machineIP string
					addrs, _ := net.InterfaceAddrs()
					for _, addr := range addrs {
						if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
							if ipNet.IP.To4() != nil {
								machineIP = ipNet.IP.String()
							}
						}
					}

					fmt.Println("Формирование запроса POST...")
					// request, err := http.NewRequest("POST", putString, bytes.NewBuffer(metricsJSON))
					request, err := http.NewRequest("POST", putString, bytes.NewBuffer([]byte(encryptedMessageBase64)))
					if err != nil {
						fmt.Println("Ошибка формирования запроса:", err)
					}
					fmt.Println("Формирование запроса POST прошло успешно")

					fmt.Println("Установка заголовков...")
					//request.Header.Set("Content-Encoding", "gzip")
					request.Header.Set("Content-Type", "application/json")
					if config.hashKey != "" {
						fmt.Println("Формирование hash...")
						hash := calcHash(metricsJSON, config.hashKey)
						fmt.Println("Формирование hash прошло успешно")
						request.Header.Set("HashSHA256", hash)
					}
					request.Header.Set("X-Real-IP", machineIP)
					fmt.Println("Установка заголовков прошла успешно")

					fmt.Println("Отправка запроса...")
					response, err := client.Do(request)
					if err != nil {
						fmt.Println("Ошибка отправки запроса:", err)
						continue
					}
					response.Body.Close()
					if response.StatusCode != http.StatusOK {
						fmt.Println("Сервер вернул ответ отличный от 200:", response.Status)
						continue
					}
					fmt.Println("Отправка запроса прошла успешно")
					fmt.Println("Reporting completed")

				case "gRPC":
					req := &pb.UpdateMetricsRequest{
						Id:   metrics.ID,
						Type: metrics.MType,
					}
					switch metrics.MType {
					case constants.GaugeType:
						req.Value = *metrics.Value
					case constants.CounterType:
						req.Delta = *metrics.Delta
					}
					_, err := gRPCMetricsClient.UpdateMetrics(ctx, req)
					if err != nil {
						fmt.Printf("ошибка при отправке метрик с помощью gRPC: %s\n", err.Error())
						continue
					}
				}

			}
		}
	}
}

// NewAgent - конструктор агента.
func NewAgent(store mapstorage.Storage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}

func getPublicKey(publicPathKey string) (*rsa.PublicKey, error) {
	file, err := os.ReadFile(publicPathKey)
	if err != nil {
		return nil, fmt.Errorf("получение publicKey: открытие файла: %s", err.Error())
	}

	block, rest := pem.Decode(file)
	if len(rest) > 0 {
		return nil, fmt.Errorf("получение publicKey: неожиданные данные после PEM-блока")
	}
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("получение publicKey: неверный формат файла")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("получение publicKey: неверный формат файла: %s", err.Error())
	}

	return publicKey, nil
}

func encryptBody(body []byte, publicKey *rsa.PublicKey) (string, error) {
	encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, body)
	if err != nil {
		return "", fmt.Errorf("шифрование тела запроса: %s", err.Error())
	}

	res := base64.StdEncoding.EncodeToString(encryptedBody)

	return res, nil
}
