package server

import (
	"flag"
	"os"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type Config struct {
	Addr            string
	GaugeType       string
	CounterType     string
	UpdateMethod    string
	GetValueMethod  string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DatabaseDNS     string
}

func MewConfig() (*Config, error) {
	var flagRunAddr, flagFileStoragePath, flagDatabaseDNS string
	var flagStoreInterval int
	var flagRestore bool

	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 300, "the interval for saving to a file")
	flag.StringVar(&flagFileStoragePath, "f", "metrics.txt", "the path to the metric file")
	flag.BoolVar(&flagRestore, "r", true, "get saved metrics from a file")
	flag.StringVar(&flagDatabaseDNS, "d", "host=localhost port=5433 user=postgres password=password dbname=metalecoll sslmode=disable", "database DNS")
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
	if envDatabaseDNS := os.Getenv("DATABASE_DSN"); envDatabaseDNS != "" {
		flagDatabaseDNS = envDatabaseDNS
	}

	return &Config{
		Addr:            flagRunAddr,
		GaugeType:       constants.GaugeType,
		CounterType:     constants.CounterType,
		UpdateMethod:    constants.UpdateMethod,
		GetValueMethod:  constants.GetValueMethod,
		LogLevel:        "debug",
		StoreInterval:   flagStoreInterval,
		FileStoragePath: flagFileStoragePath,
		Restore:         flagRestore,
		DatabaseDNS:     flagDatabaseDNS,
	}, nil
}
