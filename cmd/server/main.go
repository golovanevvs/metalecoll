package main

import (
	"github.com/golovanevvs/metalecoll/internal/server/server"
)

func main() {
	config := server.MewConfig()
	server.Start(config)
}

//Для запуска теста iter1

//surface
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\dev\\projects\\yapracticum\\metalecoll\\cmd\\server\\server

//home
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\cmd\\server\\server
