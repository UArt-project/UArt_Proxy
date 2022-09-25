// Package server is responsible for server creation
// and handling.
package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/UArt-project/UArt-proxy/cmd/server/config"
	"github.com/UArt-project/UArt-proxy/pkg/logger"
)

// Server defines an http server.
type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

// NewServer returns new instance of server.Server with specified configuration.
func NewServer(serverConfig *config.Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              serverConfig.Address,
			Handler:           serverConfig.Handler,
			TLSConfig:         nil,
			ReadTimeout:       serverConfig.ReadTimeout,
			ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
			WriteTimeout:      serverConfig.WriteTimeout,
			IdleTimeout:       serverConfig.IdleTimeout,
			MaxHeaderBytes:    0,
			TLSNextProto:      nil,
			ConnState:         nil,
			ErrorLog:          serverConfig.ErrorLog,
			BaseContext:       nil,
			ConnContext:       nil,
		},
		logger: serverConfig.ServerLogger,
	}
}

// StartListening runs server.Server.
func (s *Server) StartListening(stopChan chan struct{}) {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("calling ListenAndServe resulted in %v", err)
		}

		close(stopChan)
	}()

	s.logger.Info("listening on address %s", s.httpServer.Addr)
}

// Shutdown gracefully stops server.Server.
func (s *Server) Shutdown() {
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		s.logger.Error("server shutdown resulted in %v", err)
	}
}
