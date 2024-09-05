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
}

func MewConfig() (*Config, error) {
	var flagRunAddr, flagFileStoragePath string
	var flagStoreInterval int
	var flagRestore bool

	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.IntVar(&flagStoreInterval, "i", 15, "the interval for saving to a file")
	flag.StringVar(&flagFileStoragePath, "f", "storage/filestorage/metrics.txt", "the path to the metric file")
	flag.BoolVar(&flagRestore, "r", true, "get saved metrics from a file")
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
	}, nil
}
