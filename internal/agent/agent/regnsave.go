package agent

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/agent/model"
)

func regNSave(ag *agent) {
	const numJobs = 32
	var w int

	w = 3

	fmt.Println("Запуск regNSave")

	jobs := make(chan int, numJobs)
	results := make(chan model.Metric, numJobs)

	for i := 1; i <= w; i++ {
		go regRTMetWorker(jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	for a := 1; a <= numJobs; a++ {
		ag.store.SaveMetric(<-results)
	}

	close(jobs)
	fmt.Println(ag.store.GetMetricsMap())
	fmt.Println("Завершение regNSave")
}
