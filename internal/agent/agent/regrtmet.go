package agent

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

var pCount int64

func RegisterMetrics() {
	var rtMet runtime.MemStats
	var newMet model.Metric
	var rV float64

	pCount++

	rV = rand.Float64()

	fmt.Println("Updating №", pCount)

	// Сбор метрик и сохранение метрик в map-хранилище
	runtime.ReadMemStats(&rtMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Alloc",
		MetValue: float64(rtMet.Alloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "BuckHashSys",
		MetValue: float64(rtMet.BuckHashSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Frees",
		MetValue: float64(rtMet.Frees),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCCPUFraction",
		MetValue: float64(rtMet.GCCPUFraction),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCSys",
		MetValue: float64(rtMet.GCSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapAlloc",
		MetValue: float64(rtMet.HeapAlloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapIdle",
		MetValue: float64(rtMet.HeapIdle),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapInuse",
		MetValue: float64(rtMet.HeapInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapObjects",
		MetValue: float64(rtMet.HeapObjects),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapReleased",
		MetValue: float64(rtMet.HeapReleased),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapSys",
		MetValue: float64(rtMet.HeapSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "LastGC",
		MetValue: float64(rtMet.LastGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheInuse",
		MetValue: float64(rtMet.MCacheInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheSys",
		MetValue: float64(rtMet.MCacheSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanInuse",
		MetValue: float64(rtMet.MSpanInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanSys",
		MetValue: float64(rtMet.MSpanSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Mallocs",
		MetValue: float64(rtMet.Mallocs),
	}
	ag.store.SaveMetric(newMet)
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NextGC",
		MetValue: float64(rtMet.NextGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumForcedGC",
		MetValue: float64(rtMet.NumForcedGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumGC",
		MetValue: float64(rtMet.NumGC),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "OtherSys",
		MetValue: float64(rtMet.OtherSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "PauseTotalNs",
		MetValue: float64(rtMet.PauseTotalNs),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackInuse",
		MetValue: float64(rtMet.StackInuse),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackSys",
		MetValue: float64(rtMet.StackSys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Sys",
		MetValue: float64(rtMet.Sys),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "TotalAlloc",
		MetValue: float64(rtMet.TotalAlloc),
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.CounterType,
		MetName:  "PollCount",
		MetValue: pCount,
	}
	ag.store.SaveMetric(newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "RandomValue",
		MetValue: rV,
	}
	ag.store.SaveMetric(newMet)
}
