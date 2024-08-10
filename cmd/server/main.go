package main

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/handlers"
)

func main() {
	fmt.Println("Запущен сервер:", constants.Addr)
	err := http.ListenAndServe(constants.Addr, http.HandlerFunc(handlers.MainHandler))
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
}

//Для запуска теста iter1

//surface
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\dev\\projects\\yapracticum\\metalecoll\\cmd\\server\\server

//home
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\cmd\\server\\server
