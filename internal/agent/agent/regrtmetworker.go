package agent

import (
	"fmt"
	"math/rand/v2"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

func regRTMetWorker(r reg, jobs <-chan int, results chan<- model.Metric) {
	for j := range jobs {
		switch j {
		case 0:
			results <- model.Metric{Type: constants.GaugeType, Name: "Alloc", Value: float64(r.rtMet.Alloc)}
		case 1:
			results <- model.Metric{Type: constants.GaugeType, Name: "BuckHashSys", Value: float64(r.rtMet.BuckHashSys)}
		case 2:
			results <- model.Metric{Type: constants.GaugeType, Name: "Frees", Value: float64(r.rtMet.Frees)}
		case 3:
			results <- model.Metric{Type: constants.GaugeType, Name: "GCCPUFraction", Value: float64(r.rtMet.GCCPUFraction)}
		case 4:
			results <- model.Metric{Type: constants.GaugeType, Name: "GCSys", Value: float64(r.rtMet.GCSys)}
		case 5:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapAlloc", Value: float64(r.rtMet.HeapAlloc)}
		case 6:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapIdle", Value: float64(r.rtMet.HeapIdle)}
		case 7:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapInuse", Value: float64(r.rtMet.HeapInuse)}
		case 8:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapObjects", Value: float64(r.rtMet.HeapObjects)}
		case 9:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapReleased", Value: float64(r.rtMet.HeapReleased)}
		case 10:
			results <- model.Metric{Type: constants.GaugeType, Name: "HeapSys", Value: float64(r.rtMet.HeapSys)}
		case 11:
			results <- model.Metric{Type: constants.GaugeType, Name: "LastGC", Value: float64(r.rtMet.LastGC)}
		case 12:
			results <- model.Metric{Type: constants.GaugeType, Name: "Lookups", Value: float64(r.rtMet.Lookups)}
		case 13:
			results <- model.Metric{Type: constants.GaugeType, Name: "MCacheInuse", Value: float64(r.rtMet.MCacheInuse)}
		case 14:
			results <- model.Metric{Type: constants.GaugeType, Name: "MCacheSys", Value: float64(r.rtMet.MCacheSys)}
		case 15:
			results <- model.Metric{Type: constants.GaugeType, Name: "MSpanInuse", Value: float64(r.rtMet.MSpanInuse)}
		case 16:
			results <- model.Metric{Type: constants.GaugeType, Name: "MSpanSys", Value: float64(r.rtMet.MSpanSys)}
		case 17:
			results <- model.Metric{Type: constants.GaugeType, Name: "Mallocs", Value: float64(r.rtMet.Mallocs)}
		case 18:
			results <- model.Metric{Type: constants.GaugeType, Name: "NextGC", Value: float64(r.rtMet.NextGC)}
		case 19:
			results <- model.Metric{Type: constants.GaugeType, Name: "NumForcedGC", Value: float64(r.rtMet.NumForcedGC)}
		case 20:
			results <- model.Metric{Type: constants.GaugeType, Name: "NumGC", Value: float64(r.rtMet.NumGC)}
		case 21:
			results <- model.Metric{Type: constants.GaugeType, Name: "OtherSys", Value: float64(r.rtMet.OtherSys)}
		case 22:
			results <- model.Metric{Type: constants.GaugeType, Name: "PauseTotalNs", Value: float64(r.rtMet.PauseTotalNs)}
		case 23:
			results <- model.Metric{Type: constants.GaugeType, Name: "StackInuse", Value: float64(r.rtMet.StackInuse)}
		case 24:
			results <- model.Metric{Type: constants.GaugeType, Name: "StackSys", Value: float64(r.rtMet.StackSys)}
		case 25:
			results <- model.Metric{Type: constants.GaugeType, Name: "Sys", Value: float64(r.rtMet.Sys)}
		case 26:
			results <- model.Metric{Type: constants.GaugeType, Name: "TotalAlloc", Value: float64(r.rtMet.TotalAlloc)}
		case 27:
			results <- model.Metric{Type: constants.GaugeType, Name: "RandomValue", Value: rand.Float64()}
		case 28:
			pCount++
			results <- model.Metric{Type: constants.CounterType, Name: "PollCount", Value: pCount}
		case 29:
			results <- model.Metric{Type: constants.GaugeType, Name: "TotalMemory", Value: float64(r.vm.Total)}
		case 30:
			results <- model.Metric{Type: constants.GaugeType, Name: "FreeMemory", Value: float64(r.vm.Free)}
		case 31:
			for i, c := range r.cpuUt {
				results <- model.Metric{Type: constants.GaugeType, Name: fmt.Sprintf("CPUutilization%v", i+1), Value: c}
			}
		}
	}
}
