package server

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	tgauge   = "gauge"
	tcounter = "counter"
)

type Metric struct {
	Name    string
	Gauge   float64
	Counter int64
}

type MemStorage struct {
	metrics map[string]Metric
}

type Storage interface {
	Add(metric Metric) error
}

type mainHandler struct{}

func Start(config *Config) error {
	return http.ListenAndServe(config.bindAddr, nil)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sbody := strings.Split(string(body), "/")
	if len(sbody) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mT := sbody[2]
	//mN := sbody[3]
	mV := sbody[4]

	if mT == tgauge {
		mVgauge, err := strconv.ParseFloat(mV, 64)
		if err != nil || mVgauge < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if mT == tcounter {
		mVcounter, err := strconv.ParseInt(mV, 10, 64)
		if err != nil || mVcounter < 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
