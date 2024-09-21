package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type Config struct {
	Server          Server
	Logger          Logger
	Storage         Storage
	MetricTypeNames MetricTypeNames
}

type Server struct {
	Addr          string
	StoreInterval int
}

type Storage struct {
	Restore         bool
	FileStoragePath string
	DatabaseDSN     string
}

type MetricTypeNames struct {
	GaugeType   string
	CounterType string
}

type Logger struct {
	LogLevel string
}

func Mew() (*Config, error) {
	var flagRunAddr, flagFileStoragePath, flagDatabaseDSN string
	var flagStoreInterval int
	var flagRestore bool

	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 300, "the interval for saving to a file")
	// flag.IntVar(&flagStoreInterval, "i", 15, "the interval for saving to a file")
	flag.StringVar(&flagFileStoragePath, "f", "metrics.txt", "the path to the metric file")
	flag.BoolVar(&flagRestore, "r", true, "get saved metrics from a file")
	flag.StringVar(&flagDatabaseDSN, "d", "", "database DSN")
	// flag.StringVar(&flagDatabaseDSN, "d", "host=localhost port=5433 user=postgres password=password dbname=metalecoll sslmode=disable", "database DSN")
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

	return &Config{
		Server{
			Addr:          flagRunAddr,
			StoreInterval: flagStoreInterval,
		},
		Logger{
			LogLevel: "debug",
		},
		Storage{
			Restore:         flagRestore,
			FileStoragePath: flagFileStoragePath,
			DatabaseDSN:     flagDatabaseDSN,
		},
		MetricTypeNames{
			GaugeType:   constants.GaugeType,
			CounterType: constants.CounterType,
		},
	}, nil
}
