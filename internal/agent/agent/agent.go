package agent

import (
	"flag"
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
var flagRunAddr string
var flagRepInt int
var flagPollInt int

func Start() {
	var pollInterval int
	var reportInterval int
	var t1, t2 int
	var putString string

	parseFlags()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	pollInterval = flagPollInt
	reportInterval = flagRepInt

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
			putString = fmt.Sprintf("http://%s/update/%s/%s/%v", flagRunAddr, value.MetType, value.MetName, value.MetValue)
			request, err := client.Post(putString, constants.AContentType, nil)
			if err != nil {
				fmt.Println("Ошибка отправки POST-запроса:", err)
			}
			defer request.Body.Close()
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

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", constants.Addr, "address and port of server")
	flag.IntVar(&flagRepInt, "r", 10, "reportInterval")
	flag.IntVar(&flagPollInt, "p", 2, "pollInterval")
	flag.Parse()
}
