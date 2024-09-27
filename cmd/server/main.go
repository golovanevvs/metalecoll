package main

import (
	"log"
	"time"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/server"
)

func main() {
	c, err := config.New()
	if err != nil {
		time.Sleep(10 * time.Second)
		log.Fatalf("Ошибка конфигурирования сервера: %v", err)
	}
	err = server.Start(c)
	if err != nil {
		time.Sleep(10 * time.Second)
		log.Fatalf("Ошибка: %v", err)
	}
	time.Sleep(5 * time.Second)
}
