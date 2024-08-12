package agent

import (
	"time"

	"github.com/golovanevvs/metalecoll/internal/agent/storage/amapstorage"
)

type agent struct {
	store        amapstorage.AStorage
	updatingTime int
	sendingTime  int
}

var ag *agent

func Start() {
	var updatingTime int = 2
	var sendingTime int = 10
	var t1, t2 int

	store := amapstorage.NewStorage()

	ag = NewAgent(store, updatingTime, sendingTime)

	t1 = ag.sendingTime / ag.updatingTime
	t2 = ag.sendingTime % ag.updatingTime

	for {
		for i := 0; i < t1; i++ {
			RegisterMetrics()
			time.Sleep(time.Duration(ag.updatingTime) * time.Second)
		}
		time.Sleep(time.Duration(t2) * time.Second)
		//TODO: добавить отправку метрик методом POST
	}
}

func NewAgent(store amapstorage.AStorage, updatingTime, sendingTime int) *agent {
	s := &agent{
		store:        store,
		updatingTime: updatingTime,
		sendingTime:  sendingTime,
	}
	return s
}
