package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/UArt-project/UArt-proxy/api/v1/rest"
	"github.com/UArt-project/UArt-proxy/cmd/server"
	"github.com/UArt-project/UArt-proxy/cmd/server/config"
	"github.com/UArt-project/UArt-proxy/internal/service"
	"github.com/UArt-project/UArt-proxy/pkg/cache"
	"github.com/UArt-project/UArt-proxy/pkg/clients/marketclient"
	"github.com/UArt-project/UArt-proxy/pkg/configreader"
	"github.com/UArt-project/UArt-proxy/pkg/cors"
	"github.com/UArt-project/UArt-proxy/pkg/logger"
	"github.com/UArt-project/UArt-proxy/pkg/workerpool"
)

const configFile = "config.yaml"

func main() {
	mainLogger := logger.NewLogger(os.Stdout, "main")

	mainLogger.Info("Starting the application...")

	err := configreader.SetConfigFile(configFile)
	if err != nil {
		mainLogger.Fatal("setting the config file: %v", err)
	}

	marketClient := marketclient.NewMarketServiceClient(configreader.GetString("market_service_url"))
	pool := workerpool.NewPool(configreader.GetInt("worker_pool_size"))
	appCache := cache.NewLocalCache(15 * time.Second)
	appService := service.NewService(marketClient, pool, appCache)
	restLogger := logger.NewLogger(os.Stdout, "rest")
	restAPI := rest.NewRESTApi(appService, restLogger)
	serverLogger := logger.NewLogger(os.Stdout, "server")
	serverConfig := getServerConfig(cors.EnableCORS(restAPI), nil, serverLogger)
	restServer := server.NewServer(serverConfig)
	serverStopChan := make(chan struct{})

	restServer.StartListening(serverStopChan)

	serverWG := new(sync.WaitGroup)
	numberOfServersRunning := 1

	serverWG.Add(numberOfServersRunning)

	go func(wg *sync.WaitGroup) {
		<-serverStopChan

		wg.Done()
	}(serverWG)

	serverWG.Wait()
}

// getServerConfig reads the server configuration from the config file.
func getServerConfig(handler http.Handler, errorLog *log.Logger, serverLogger *logger.Logger) *config.Config {
	var (
		address          = configreader.GetString("server.address")
		readTime         = configreader.GetDuration("server.readTime")
		writeTime        = configreader.GetDuration("server.writeTime")
		idleTime         = configreader.GetDuration("server.idleTime")
		readerHeaderTime = configreader.GetDuration("server.readerHeaderTime")
	)

	return &config.Config{
		Address:           address,
		ReadTimeout:       readTime,
		WriteTimeout:      writeTime,
		IdleTimeout:       idleTime,
		ReadHeaderTimeout: readerHeaderTime,
		ErrorLog:          errorLog,
		ServerLogger:      serverLogger,
		Handler:           handler,
	}
}
