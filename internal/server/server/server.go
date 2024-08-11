package server

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

type server struct {
	store mapstorage.Storage
}

func Start() {
	store := mapstorage.NewStorage()
	srv := NewServer(store)
	fmt.Println("Запущен сервер:", constants.Addr)
	err := http.ListenAndServe(constants.Addr, srv)
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
}

func NewServer(store mapstorage.Storage) *server {
	s := &server{
		store: store,
	}
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	MainHandler(w, r)
}
