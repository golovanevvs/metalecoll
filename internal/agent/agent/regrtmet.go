package agent

import (
	"runtime"

	"github.com/golovanevvs/metalecoll/internal/agent/autil"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
)

func RegisterMetrics() {
	var rtMet runtime.MemStats
	var newMet model.Metric
	newMet = model.Metric{
		MetType:  constants.CounterType,
		MetName:  "PollCount",
		MetValue: 1,
	}

	runtime.ReadMemStats(&rtMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Alloc",
		MetValue: float64(rtMet.Alloc),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "BuckHashSys",
		MetValue: float64(rtMet.BuckHashSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Frees",
		MetValue: float64(rtMet.Frees),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCCPUFraction",
		MetValue: float64(rtMet.GCCPUFraction),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "GCSys",
		MetValue: float64(rtMet.GCSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapAlloc",
		MetValue: float64(rtMet.HeapAlloc),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapIdle",
		MetValue: float64(rtMet.HeapIdle),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapInuse",
		MetValue: float64(rtMet.HeapInuse),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapObjects",
		MetValue: float64(rtMet.HeapObjects),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapReleased",
		MetValue: float64(rtMet.HeapReleased),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "HeapSys",
		MetValue: float64(rtMet.HeapSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "LastGC",
		MetValue: float64(rtMet.LastGC),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheInuse",
		MetValue: float64(rtMet.MCacheInuse),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MCacheSys",
		MetValue: float64(rtMet.MCacheSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanInuse",
		MetValue: float64(rtMet.MSpanInuse),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "MSpanSys",
		MetValue: float64(rtMet.MSpanSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Mallocs",
		MetValue: float64(rtMet.Mallocs),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NextGC",
		MetValue: float64(rtMet.NextGC),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Lookups",
		MetValue: float64(rtMet.Lookups),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumForcedGC",
		MetValue: float64(rtMet.NumForcedGC),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "NumGC",
		MetValue: float64(rtMet.NumGC),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "OtherSys",
		MetValue: float64(rtMet.OtherSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "PauseTotalNs",
		MetValue: float64(rtMet.PauseTotalNs),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackInuse",
		MetValue: float64(rtMet.StackInuse),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "StackSys",
		MetValue: float64(rtMet.StackSys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "Sys",
		MetValue: float64(rtMet.Sys),
	}
	autil.SM(ag.store, newMet)
	newMet = model.Metric{
		MetType:  constants.GaugeType,
		MetName:  "TotalAlloc",
		MetValue: float64(rtMet.TotalAlloc),
	}
	autil.SM(ag.store, newMet)
}
