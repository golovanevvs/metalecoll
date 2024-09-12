package agent

import (
	"fmt"
	"math/rand"
	"runtime"

	"github.com/golovanevvs/metalecoll/internal/agent/storage"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

var pCount int64

func RegisterMetrics() {
	var rtMet runtime.MemStats
	var newMet model.Metric
	var rV float64

	pCount++

	rV = rand.Float64()

	fmt.Println("Updating №", pCount)

	runtime.ReadMemStats(&rtMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Alloc",
		MetValue: float64(rtMet.Alloc),
	}
	storage.SM(ag.store, newMet)

	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "BuckHashSys",
		MetValue: float64(rtMet.BuckHashSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Frees",
		MetValue: float64(rtMet.Frees),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCCPUFraction",
		MetValue: float64(rtMet.GCCPUFraction),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCSys",
		MetValue: float64(rtMet.GCSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapAlloc",
		MetValue: float64(rtMet.HeapAlloc),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapIdle",
		MetValue: float64(rtMet.HeapIdle),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapInuse",
		MetValue: float64(rtMet.HeapInuse),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapObjects",
		MetValue: float64(rtMet.HeapObjects),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapReleased",
		MetValue: float64(rtMet.HeapReleased),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapSys",
		MetValue: float64(rtMet.HeapSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "LastGC",
		MetValue: float64(rtMet.LastGC),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheInuse",
		MetValue: float64(rtMet.MCacheInuse),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheSys",
		MetValue: float64(rtMet.MCacheSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanInuse",
		MetValue: float64(rtMet.MSpanInuse),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanSys",
		MetValue: float64(rtMet.MSpanSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Mallocs",
		MetValue: float64(rtMet.Mallocs),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NextGC",
		MetValue: float64(rtMet.NextGC),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumForcedGC",
		MetValue: float64(rtMet.NumForcedGC),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumGC",
		MetValue: float64(rtMet.NumGC),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "OtherSys",
		MetValue: float64(rtMet.OtherSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "PauseTotalNs",
		MetValue: float64(rtMet.PauseTotalNs),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackInuse",
		MetValue: float64(rtMet.StackInuse),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackSys",
		MetValue: float64(rtMet.StackSys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Sys",
		MetValue: float64(rtMet.Sys),
	}
	storage.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "TotalAlloc",
		MetValue: float64(rtMet.TotalAlloc),
	}
	storage.SM(ag.store, newMet)

	newMet = model.Metric{
		MetType:  constants.CounterType,
		MetName:  "PollCount",
		MetValue: pCount,
	}
	storage.SM(ag.store, newMet)

	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "RandomValue",
		MetValue: rV,
	}
	storage.SM(ag.store, newMet)

}
