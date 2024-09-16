package main

import (
	"log"

	"github.com/golovanevvs/metalecoll/internal/server/server"
)

func main() {
	config, err := server.MewConfig()
	if err != nil {
		log.Fatalf("Ошибка конфигурирования сервера: %v", err)
	}
	server.Start(config)
	// if err != nil {
	// 	log.Fatalf("Ошибка: %v", err)
	// }
}

//Для запуска теста iter1

//surface
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\dev\\projects\\yapracticum\\metalecoll\\cmd\\server\\server

//home
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\cmd\\server\\server
