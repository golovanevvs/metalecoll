package agent

import "fmt"

func sendMetrics(metrics [][]Metrics, urlString string, hashKey string, limit int) {
	limit = 1

	fmt.Printf("len(metrics): %d\n", len(metrics))

	// создание буферизованного канала для принятия задач в воркер
	jobs := make(chan []Metrics, len(metrics))

	// создание буферизованного канала для отправки результатов
	results := make(chan string, len(metrics))

	// создание и запуск воркеров
	for i := 0; i < limit; i++ {
		go sendMetWorker(i, urlString, hashKey, jobs, results)
	}

	// отправка задачи в канал задач
	for j := 0; j < len(metrics); j++ {
		jobs <- metrics[j]
	}

	//получение результатов из канала результатов
	for a := 0; a < len(metrics); a++ {
		fmt.Println(<-results)
	}

	//закрытие канала на стороне отправителя
	fmt.Println("Закрытие канала jobs")
	close(jobs)
}
