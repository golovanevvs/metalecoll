package agent

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

var pCount int64

func RegisterMetrics(ag *agent) {
	var rtMet runtime.MemStats
	var newMet model.Metric
	var rV float64

	pCount++

	rV = rand.Float64()

	fmt.Println("Updating №", pCount)

	// Сбор метрик и сохранение метрик в map-хранилище
	runtime.ReadMemStats(&rtMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Alloc",
		Value: float64(rtMet.Alloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "BuckHashSys",
		Value: float64(rtMet.BuckHashSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Frees",
		Value: float64(rtMet.Frees),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "GCCPUFraction",
		Value: float64(rtMet.GCCPUFraction),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "GCSys",
		Value: float64(rtMet.GCSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapAlloc",
		Value: float64(rtMet.HeapAlloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapIdle",
		Value: float64(rtMet.HeapIdle),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapInuse",
		Value: float64(rtMet.HeapInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapObjects",
		Value: float64(rtMet.HeapObjects),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapReleased",
		Value: float64(rtMet.HeapReleased),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "HeapSys",
		Value: float64(rtMet.HeapSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "LastGC",
		Value: float64(rtMet.LastGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Lookups",
		Value: float64(rtMet.Lookups),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "MCacheInuse",
		Value: float64(rtMet.MCacheInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "MCacheSys",
		Value: float64(rtMet.MCacheSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "MSpanInuse",
		Value: float64(rtMet.MSpanInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "MSpanSys",
		Value: float64(rtMet.MSpanSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Mallocs",
		Value: float64(rtMet.Mallocs),
	}
	ag.store.SaveMetric(newMet)
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "NextGC",
		Value: float64(rtMet.NextGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Lookups",
		Value: float64(rtMet.Lookups),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "NumForcedGC",
		Value: float64(rtMet.NumForcedGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "NumGC",
		Value: float64(rtMet.NumGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "OtherSys",
		Value: float64(rtMet.OtherSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "PauseTotalNs",
		Value: float64(rtMet.PauseTotalNs),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "StackInuse",
		Value: float64(rtMet.StackInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "StackSys",
		Value: float64(rtMet.StackSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "Sys",
		Value: float64(rtMet.Sys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "TotalAlloc",
		Value: float64(rtMet.TotalAlloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.CounterType,
		Name:  "PollCount",
		Value: pCount,
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		Type:  constants.GaugeType,
		Name:  "RandomValue",
		Value: rV,
	}
	ag.store.SaveMetric(newMet)
}
