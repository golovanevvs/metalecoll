package main

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/handlers"
	"github.com/golovanevvs/metalecoll/internal/server/server"
)

func main() {
	fmt.Println("Запущен сервер:", server.Addr)
	err := http.ListenAndServe(server.Addr, http.HandlerFunc(handlers.MainHandler))
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
}

//TODO Обернуть запрос в структуру
//TODO Доработать архитектуру

//Для запуска теста iter1

//surface
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\dev\\projects\\yapracticum\\metalecoll\\cmd\\server\\server

//home
//metricstest-windows-amd64 -test.v -test.run=^TestIteration1$ -binary-path=C:\\Golovanev\\Dev\\Projects\\YaPracticum\\metalecoll\\cmd\\server
