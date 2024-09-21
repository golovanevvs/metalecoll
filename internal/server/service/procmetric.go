package service

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/model"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

func ProcMetric(recMet model.Metric, store mapstorage.Storage) *model.Metric {
	fmt.Printf("____________________Запущен ProcMetric...\n")

	var newValue any

	switch recMet.MetType {

	case constants.GaugeType:
		fmt.Printf("%v\n", constants.GaugeType)
		newValue = recMet.MetValue.(float64)
		fmt.Printf("newValue = %v\n", newValue)

	case constants.CounterType:
		fmt.Printf("%v\n", constants.CounterType)
		if getValue, err := store.GetMetric(recMet.MetName); err != nil {
			fmt.Printf("Нет в мапе\n")
			newValue = recMet.MetValue.(int64)
			fmt.Printf("newValue = %v\n", newValue)
		} else {
			fmt.Printf("Есть в мапе\n")
			fmt.Printf("recMet.MetValue = %v\n", recMet.MetValue.(int64))
			fmt.Printf("getValue.MetValue = %v\n", getValue.MetValue.(int64))
			newValue = (getValue.MetValue.(int64)) + recMet.MetValue.(int64)
			fmt.Printf("newValue = %v\n", newValue)
		}
	}

	fmt.Printf("____________________Завершён ProcMetric\n")
	return &model.Metric{
		MetType:  recMet.MetType,
		MetName:  recMet.MetName,
		MetValue: newValue,
	}
}
