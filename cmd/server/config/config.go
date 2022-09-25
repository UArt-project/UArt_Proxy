// Package config is used to store all the needed options to run a server.
package config

import (
	"log"
	"net/http"
	"time"

	"github.com/UArt-project/UArt-proxy/pkg/logger"
)

// Config consists of data needed for server configuration.
type Config struct {
	Address           string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	ErrorLog          *log.Logger
	ServerLogger      *logger.Logger
	Handler           http.Handler
}
