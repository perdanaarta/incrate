package api

import (
	"fmt"
	"net/http"
)

type ServerAPI struct {
	Host string
	Port int
}

func NewAPIsServer(host string, port int) *ServerAPI {
	return &ServerAPI{
		Host: host,
		Port: port,
	}
}

func (s *ServerAPI) Run() error {
	router := NewAPIsRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
