package agent

import (
	"flag"
	"os"
	"strconv"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type config struct {
	addr           string
	gaugeType      string
	counterType    string
	updateMethod   string
	getValueMethod string
	pollInterval   int
	reportInterval int
	hashKey        string
	rateLimit      int
}

// NewConfig - конструктор конфигурации агента.
func NewConfig() *config {
	var flagRunAddr string
	var flagRepInt int
	var flagPollInt int
	var flagHashKey string
	var flagRateLimit int

	flag.StringVar(&flagRunAddr, "a", constants.AddrA, "address and port of server")
	flag.IntVar(&flagRepInt, "r", 10, "reportInterval")
	flag.IntVar(&flagPollInt, "p", 2, "pollInterval")
	flag.StringVar(&flagHashKey, "k", "", "hash key")
	flag.IntVar(&flagRateLimit, "l", 3, "rate limit")
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
	if flagRateLimit == 0 {
		flagRateLimit = 1
	}

	return &config{
		addr:           flagRunAddr,
		gaugeType:      constants.GaugeType,
		counterType:    constants.CounterType,
		updateMethod:   constants.UpdateMethod,
		getValueMethod: constants.GetValueMethod,
		pollInterval:   flagPollInt,
		reportInterval: flagRepInt,
		hashKey:        flagHashKey,
		rateLimit:      flagRateLimit,
	}
}
