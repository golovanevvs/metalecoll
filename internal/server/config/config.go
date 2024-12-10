// Модуль config предназначен для определения конфигурации приложения.
package config

import (
	"flag"
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
	Addr          string
	StoreInterval int
}

// Storage - структура конфигурации хранилища.
type Storage struct {
	Restore         bool
	FileStoragePath string
	DatabaseDSN     string
}

// Logger - структура конфигурации логгера.
type Logger struct {
	LogLevel logrus.Level
}

// Crypto - структура конфигурации шифрования.
type Crypto struct {
	HashKey string
}

// NewConfig - конструктор конфигурации.
func NewConfig() (*Config, error) {
	var flagRunAddr, flagFileStoragePath, flagDatabaseDSN, flagHashKey string
	var flagStoreInterval int
	var flagRestore bool

	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 15, "the interval for saving to a file")
	flag.StringVar(&flagFileStoragePath, "f", "metrics.txt", "the path to the metric file")
	flag.BoolVar(&flagRestore, "r", true, "get saved metrics from a file")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database DSN")
	//flag.StringVar(&flagDatabaseDSN, "d", "host=localhost port=5433 user=postgres password=password dbname=metalecoll sslmode=disable", "database DSN")
	flag.StringVar(&flagHashKey, "k", "", "hash key")
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

	return &Config{
		Server{
			Addr:          flagRunAddr,       // флаг: адрес сервера
			StoreInterval: flagStoreInterval, // флаг: интервал сохранения данных
		},
		Storage{
			Restore:         flagRestore,         // флаг: восстановление данных при запусае сервера
			FileStoragePath: flagFileStoragePath, // флаг: путь к файлу сохранения данных
			DatabaseDSN:     flagDatabaseDSN,     //флаг: DSN базы данных
		},
		Logger{
			logrus.DebugLevel,
		},
		Crypto{
			HashKey: flagHashKey, // флаг: хэш ключ
		},
	}, nil
}
