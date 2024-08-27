package server

import (
	"flag"
	"os"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type Config struct {
	Addr           string
	GaugeType      string
	CounterType    string
	UpdateMethod   string
	GetValueMethod string
	LogLevel       string
}

func MewConfig() *Config {
	var flagRunAddr string
	flag.StringVar(&flagRunAddr, "a", constants.AddrS, "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	return &Config{
		Addr:           flagRunAddr,
		GaugeType:      constants.GaugeType,
		CounterType:    constants.CounterType,
		UpdateMethod:   constants.UpdateMethod,
		GetValueMethod: constants.GetValueMethod,
		LogLevel:       "info",
	}
}
