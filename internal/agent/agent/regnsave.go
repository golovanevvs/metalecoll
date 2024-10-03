package agent

import (
	"fmt"
	"runtime"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type reg struct {
	rtMet runtime.MemStats
	vm    *mem.VirtualMemoryStat
	cpuUt []float64
}

// Реализация Worker Pool
func regNSave(ag *agent) {
	const numJobs = 32 // количество работ

	var w int
	var r reg
	var err error

	w = 3 // количество воркеров

	// заполнение полей объекта r
	runtime.ReadMemStats(&r.rtMet)

	r.vm, err = mem.VirtualMemory()
	if err != nil {
		r.vm = &mem.VirtualMemoryStat{}
	}

	r.cpuUt, err = cpu.Percent(0, true)
	if err != nil {
		r.cpuUt = []float64{0}
	}

	// количество работ для вывода результата
	newNumJobs := numJobs + len(r.cpuUt) - 1

	// создаём буферизованный канал для принятия задач в воркер
	jobs := make(chan int, numJobs)

	// создаём буферизованный канал для отправки результатов
	results := make(chan model.Metric, newNumJobs)

	// создаём и запускаем воркеры
	for i := 0; i < w; i++ {
		go regRTMetWorker(r, jobs, results)
	}

	// в канал задач отправляем id задачи
	for j := 0; j < numJobs; j++ {
		jobs <- j
	}

	// забираем из канала результатов результаты
	for a := 0; a < newNumJobs; a++ {
		ag.store.SaveMetric(<-results)
	}

	// закрываем канал на стороне отправителя
	close(jobs)
	close(results)

	// выводим полученное map-хранилище
	count := 0
	mapa, _ := ag.store.GetMetricsMap()
	for _, m := range mapa {
		count++
		fmt.Printf("%v. Name: %v, Value:%v\n", count, m.Name, m.Value)
	}
}
