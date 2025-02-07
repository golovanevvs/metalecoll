package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type config struct {
	Addr           string `json:"address"`
	gaugeType      string
	counterType    string
	updateMethod   string
	getValueMethod string
	PollInterval   int `json:"poll_interval"`
	ReportInterval int `json:"report_interval"`
	hashKey        string
	rateLimit      int
	PublicKeyPath  string `json:"crypto_key"`
}

// NewConfig - конструктор конфигурации агента.
func NewConfig() (*config, error) {
	var flagRunAddr string
	var flagRepInt int
	var flagPollInt int
	var flagHashKey string
	var flagRateLimit int
	var flagPublicKeyPath string
	var jsonConfigPath string
	var cfg config

	flag.StringVar(&jsonConfigPath, "config", "", "JSON config file path")

	flag.Parse()

	if envJSONConfigPath := os.Getenv("CONFIG"); envJSONConfigPath != "" {
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

		if err := json.Unmarshal(body, &cfg); err != nil {
			return nil, fmt.Errorf("десериализация JSON файла конфигурации: %s", err.Error())
		}
	}

	flag.StringVar(&flagRunAddr, "a", constants.AddrA, "address and port of server")
	flag.IntVar(&flagRepInt, "r", 10, "reportInterval")
	flag.IntVar(&flagPollInt, "p", 2, "pollInterval")
	flag.StringVar(&flagHashKey, "k", "", "hash key")
	flag.IntVar(&flagRateLimit, "l", 3, "rate limit")
	flag.StringVar(&flagPublicKeyPath, "crypto-key", "../../resources/keys/public_key_pkcs1.pem", "public key path")
	// flag.StringVar(&flagPublicKeyPath, "crypto-key", "C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\resources\\keys\\public_key_pkcs1.pem", "public key path")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envRepInt := os.Getenv("REPORT_INTERVAL"); envRepInt != "" {
		flagRepInt, _ = strconv.Atoi(envRepInt)
	}
	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		flagPollInt, _ = strconv.Atoi(envPollInt)
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		flagHashKey = envHashKey
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		flagRateLimit, _ = strconv.Atoi(envRateLimit)
	}
	if envPublicKeyPath := os.Getenv("CRYPTO_KEY"); envPublicKeyPath != "" {
		flagPublicKeyPath = envPublicKeyPath
	}
	if flagRateLimit == 0 {
		flagRateLimit = 1
	}

	if jsonConfigPath != "" {
		if flagRunAddr == "" {
			flagRunAddr = cfg.Addr
		}
		if flagRepInt == 0 {
			flagRepInt = cfg.ReportInterval
		}
		if flagPollInt == 0 {
			flagPollInt = cfg.PollInterval
		}
		if flagPublicKeyPath == "" {
			flagPublicKeyPath = cfg.PublicKeyPath
		}
	}

	return &config{
		Addr:           flagRunAddr,
		gaugeType:      constants.GaugeType,
		counterType:    constants.CounterType,
		updateMethod:   constants.UpdateMethod,
		getValueMethod: constants.GetValueMethod,
		PollInterval:   flagPollInt,
		ReportInterval: flagRepInt,
		hashKey:        flagHashKey,
		rateLimit:      flagRateLimit,
		PublicKeyPath:  flagPublicKeyPath,
	}, nil
}
