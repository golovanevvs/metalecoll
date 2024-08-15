package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/autil"
	"github.com/golovanevvs/metalecoll/internal/agent/storage/amapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

type agent struct {
	store          amapstorage.AStorage
	pollInterval   int
	reportInterval int
}

var ag *agent

func Start() {
	var pollInterval int = 2
	var reportInterval int = 10
	var t1, t2 int
	var putString string

	store := amapstorage.NewStorage()

	ag = NewAgent(store, pollInterval, reportInterval)

	t1 = ag.reportInterval / ag.pollInterval
	t2 = ag.reportInterval % ag.pollInterval

	client := &http.Client{}

	for {
		for i := 0; i <= t1; i++ {
			RegisterMetrics()
			if i != t1 {
				time.Sleep(time.Duration(ag.pollInterval) * time.Second)
			}
		}
		time.Sleep(time.Duration(t2) * time.Second)

		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Reporting...")

		mapstore, err := autil.GMM(ag.store)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, value := range mapstore {
			putString = fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", value.MetType, value.MetName, value.MetValue)
			request, err := client.Post(putString, constants.ContentType, nil)
			if err != nil {
				fmt.Println("Ошибка отправки POST-запроса:", err)
			}
			fmt.Println(request.StatusCode)
		}
		time.Sleep(time.Duration(ag.pollInterval-t2) * time.Second)
	}
}

func NewAgent(store amapstorage.AStorage, pollInterval, reportInterval int) *agent {
	s := &agent{
		store:          store,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
	}
	return s
}
