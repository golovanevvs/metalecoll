package agent

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

// getMetricsSlice получает метрики из map-хранилища и формирует слайс метрик в соответствии с dto
func getMetricsSlice(ag *agent) ([]Metrics, error) {
	var metricsdto []Metrics
	var body Metrics

	fmt.Println("Получение данных из map-хранилища...")
	mapStore, err := ag.store.GetMetricsMap()
	if err != nil {
		fmt.Println("Ошибка получения данных из map-хранилища:", err)
		return nil, err
	}
	fmt.Println("Получение данных из map-хранилища прошло успешно")

	fmt.Println("Формирование среза метрик в соответствии с dto...")

	for _, value := range mapStore {
		switch value.Type {
		case constants.GaugeType:
			v, _ := value.Value.(float64)
			body = Metrics{
				ID:    value.Name,
				MType: value.Type,
				Value: &v,
			}
		case constants.CounterType:
			v, _ := value.Value.(int64)
			body = Metrics{
				ID:    value.Name,
				MType: value.Type,
				Delta: &v,
			}
		}

		metricsdto = append(metricsdto, body)
	}
	fmt.Println("Формирование среза метрик в соответствии с dto прошло успешно")
	fmt.Println(metricsdto)

	return metricsdto, nil
}

// splitMetricsSlice делит слайс метрик на несколько частей (слайсов), количество которых соответствует RATE_LIMIT
func splitMetricsSlice(metricsSlice []Metrics, limit int) [][]Metrics {
	var newSlice [][]Metrics
	if limit == 0 || limit > len(metricsSlice) {
		limit = len(metricsSlice)
	}
	fmt.Printf("Разделение слайса метрик на %v частей...\n", limit)
	var a1, a2 int
	var partSlice []Metrics
	part := len(metricsSlice) / limit
	a1 = 0
	a2 = part
	for i := 1; i <= limit; i++ {
		if i != limit {
			partSlice = metricsSlice[a1:a2]
			a1, a2 = a2, a2+part
		} else {
			partSlice = metricsSlice[a1:]
		}
		newSlice = append(newSlice, partSlice)
	}

	fmt.Printf("Разделение слайса метрик на %v частей прошло успешно\n", limit)
	for i := range newSlice {
		fmt.Printf("Длина слайса %v: %v\n", i+1, len(newSlice[i]))
		fmt.Println("Содержимое слайса:")
		fmt.Println(newSlice[i])
	}
	return newSlice
}
