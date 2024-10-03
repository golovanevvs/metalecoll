package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
)

type agent struct {
	store          mapstorage.Storage
	pollInterval   int
	reportInterval int
}

func Start(config *config) {
	store := mapstorage.NewStorage()

	ag := NewAgent(store, config.pollInterval, config.reportInterval)

	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)

	defer pollIntTime.Stop()
	defer reportIntTime.Stop()

	for {
		select {
		case <-pollIntTime.C:
			go regNSave(ag)
		case <-reportIntTime.C:
			go sendMetrics(ag, config)
		}
	}
}

func NewAgent(store mapstorage.Storage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}

func calcHash(data []byte, key string) string {
	h := sha256.New()
	h.Write(data)
	h.Write([]byte(key))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
