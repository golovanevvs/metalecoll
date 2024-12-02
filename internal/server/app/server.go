package app

import (
	"context"
	"net/http"
)

// Server - структура сервера.
type Server struct {
	httpsrv *http.Server
}

// NewServer - конструктор сервера.
func NewServer() *Server {
	return &Server{}
}

// RunServer - запуск сервера.
func (srv *Server) RunServer(runAddress string, handler http.Handler) error {
	srv.httpsrv = &http.Server{
		Addr:    runAddress,
		Handler: handler,
	}

	return srv.httpsrv.ListenAndServe()
}

// ShutdownServer - остановка сервера.
func (srv *Server) ShutdownServer(ctx context.Context) error {
	return srv.httpsrv.Shutdown(ctx)
}
