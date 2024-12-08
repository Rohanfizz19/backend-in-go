package httpserver

import (
	"backend/config"
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type Server struct {
	config *config.HttpServer
	srv    *http.Server
	mux    http.Handler
}

func NewServer(config *config.HttpServer, router Router) *Server {
	fmt.Println(context.Background(), "http server config: ", zap.String("config", config.String()))
	return &Server{
		mux:    router.Mux(),
		config: config}
}

func (s *Server) Start() {
	s.config.Port = 8080;
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.mux,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}
	s.srv = srv

	err := srv.ListenAndServe()
	if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
		fmt.Println(context.Background(), "failed to start server %v", zap.Error(err))
	}
}
