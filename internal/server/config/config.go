// Модуль config предназначен для определения конфигурации приложения.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/sirupsen/logrus"
)

// Config - структура конфигурации приложения.
type Config struct {
	Server  Server
	Storage Storage
	Logger  Logger
	Crypto  Crypto
}

// Server - структура конфигарации сервера.
type Server struct {
	Addr          string `json:"address"`
	StoreInterval int    `json:"store_interval"`
}

// Storage - структура конфигурации хранилища.
type Storage struct {
	Restore         bool   `json:"restore"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
}

// Logger - структура конфигурации логгера.
type Logger struct {
	LogLevel logrus.Level
}

// Crypto - структура конфигурации шифрования.
type Crypto struct {
	HashKey        string
	PrivateKeyPath string `json:"crypto_key"`
}

// NewConfig - конструктор конфигурации.
func NewConfig() (*Config, error) {
	var flagRunAddr, flagFileStoragePath, flagDatabaseDSN, flagHashKey, flagPrivateKeyPath, jsonConfigPath string
	var flagStoreInterval int
	var flagRestore bool
	var config Config

	flag.StringVar(&jsonConfigPath, "config", "", "JSON config file path")

	flag.Parse()

	if envJSONConfigPath := os.Getenv("SERVER_CONFIG"); envJSONConfigPath != "" {
		jsonConfigPath = envJSONConfigPath
	}

	if jsonConfigPath != "" {
		file, err := os.Open(jsonConfigPath)
		if err != nil {
			return nil, fmt.Errorf("открытие JSON файла конфигурации: %s", err.Error())
		}
		defer file.Close()

		body, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("чтение JSON файла конфигурации: %s", err.Error())
		}

		if err := json.Unmarshal(body, &config); err != nil {
			return nil, fmt.Errorf("десериализация JSON файла конфигурации: %s", err.Error())
		}
	}

	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 15, "the interval for saving to a file")
	flag.StringVar(&flagFileStoragePath, "f", "metrics.txt", "the path to the metric file")
	flag.BoolVar(&flagRestore, "r", true, "get saved metrics from a file")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database DSN")
	//flag.StringVar(&flagDatabaseDSN, "d", "host=localhost port=5433 user=postgres password=password dbname=metalecoll sslmode=disable", "database DSN")
	flag.StringVar(&flagHashKey, "k", "", "hash key")
	flag.StringVar(&flagPrivateKeyPath, "crypto-key", "../../resources/keys", "private key path")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		v, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return nil, err
		}
		flagStoreInterval = v
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		v, err := strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
		flagRestore = v
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		flagDatabaseDSN = envDatabaseDSN
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		flagHashKey = envHashKey
	}
	if envPrivateKeyPath := os.Getenv("CRYPTO_KEY"); envPrivateKeyPath != "" {
		flagPrivateKeyPath = envPrivateKeyPath
	}

	if jsonConfigPath != "" {
		if flagRunAddr == "" {
			flagRunAddr = config.Server.Addr
		}
		if flagStoreInterval == 0 {
			flagStoreInterval = config.Server.StoreInterval
		}
		if flagFileStoragePath == "" {
			flagFileStoragePath = config.Storage.FileStoragePath
		}
		if !flagRestore {
			flagRestore = config.Storage.Restore
		}
		if flagDatabaseDSN == "" {
			flagDatabaseDSN = config.Storage.DatabaseDSN
		}
		if flagPrivateKeyPath == "" {
			flagPrivateKeyPath = config.Crypto.PrivateKeyPath
		}
	}

	return &Config{
		Server{
			Addr:          flagRunAddr,       // флаг: адрес сервера
			StoreInterval: flagStoreInterval, // флаг: интервал сохранения данных
		},
		Storage{
			Restore:         flagRestore,         // флаг: восстановление данных при запусае сервера
			FileStoragePath: flagFileStoragePath, // флаг: путь к файлу сохранения данных
			DatabaseDSN:     flagDatabaseDSN,     // флаг: DSN базы данных
		},
		Logger{
			logrus.DebugLevel,
		},
		Crypto{
			HashKey:        flagHashKey,                                                 // флаг: хэш ключ
			PrivateKeyPath: fmt.Sprintf("%s/private_key_pkcs1.pem", flagPrivateKeyPath), // флаг: путь к приватному ключу
		},
	}, nil
}
