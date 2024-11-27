package agent

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type reg struct {
	rtMet runtime.MemStats
	vm    *mem.VirtualMemoryStat
	cpuUt []float64
}

// Реализация Worker Pool
func regNSave(ag *agent) {
	var err error
	var mu sync.Mutex

	// заполнение полей объекта r
	var r reg
	runtime.ReadMemStats(&r.rtMet)

	r.vm, err = mem.VirtualMemory()
	if err != nil {
		r.vm = &mem.VirtualMemoryStat{}
	}

	r.cpuUt, err = cpu.Percent(0, true)
	if err != nil {
		r.cpuUt = []float64{0}
	}

	// количество работ
	const numJobs = 32

	// количество результатов работ
	numResults := numJobs + len(r.cpuUt) - 1

	// создание буферизованного канала для принятия задач в воркер
	jobs := make(chan int, numJobs)

	// создание буферизованного канала для отправки результатов
	results := make(chan model.Metric, numResults)

	// создание и запуск воркеров
	w := 3 // количество воркеров
	for i := 0; i < w; i++ {
		go regRTMetWorker(r, jobs, results)
	}

	// отправка id задачи в канал задач
	for j := 0; j < numJobs; j++ {
		jobs <- j
	}

	// получение результатов из канала результатов и сохранение их в map-хранилище
	for a := 0; a < numResults; a++ {
		mu.Lock()
		ag.store.SaveMetric(<-results)
		mu.Unlock()
	}

	// закрытие канала на стороне отправителя
	close(jobs)

	// вывод полученного map-хранилища
	count := 0
	mapa, _ := ag.store.GetMetricsMap()
	for _, m := range mapa {
		count++
		fmt.Printf("%d. Name: %s, Value: %v\n", count, m.Name, m.Value)
	}
}
