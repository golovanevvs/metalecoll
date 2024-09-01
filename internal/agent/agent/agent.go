package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage"
	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type agent struct {
	store          storage.Storage
	pollInterval   int
	reportInterval int
}

var ag *agent

func Start(config *config) {
	var putString string

	store := mapstorage.NewStorage()

	ag = NewAgent(store, config.pollInterval, config.reportInterval)

	client := &http.Client{}

	pollIntTime := time.NewTicker(time.Duration(ag.pollInterval) * time.Second)
	reportIntTime := time.NewTicker(time.Duration(ag.reportInterval) * time.Second)

	defer pollIntTime.Stop()
	defer reportIntTime.Stop()

	for {
		select {
		case <-pollIntTime.C:
			RegisterMetrics()
		case <-reportIntTime.C:
			fmt.Println("-------------------------------------------------------------------------")
			fmt.Println("Reporting...")

			mapstore, err := storage.GMM(ag.store)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, value := range mapstore {
				putString = fmt.Sprintf("http://%s/%s/%s/%s/%v", config.addr, config.updateMethod, value.MetType, value.MetName, value.MetValue)
				request, err := client.Post(putString, constants.AContentTypeTP, nil)
				if err != nil {
					fmt.Println("Ошибка отправки POST-запроса:", err)
				}
				defer request.Body.Close()
			}
		}
	}
}

func NewAgent(store storage.Storage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}
