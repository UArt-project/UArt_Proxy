// Package service contains main logic of the application.
package service

import (
	"fmt"

	"github.com/UArt-project/UArt-proxy/domain/marketdomain"
	"github.com/UArt-project/UArt-proxy/pkg/clients/marketclient"
	"github.com/UArt-project/UArt-proxy/pkg/workerpool"
)

// AppService provides information of main application service functionality.
type AppService interface {
	// GetMarketPage returns a page of market items.
	GetMarketPage(page int) ([]*marketdomain.MarketItem, error)
}

// Service is a main application logic.
type Service struct {
	// The market client.
	marketClient marketclient.MarketClient
	// Worker pool.
	workerPool *workerpool.WorkerPool
}

// NewService creates a new instance of the Service.
func NewService(marketClient marketclient.MarketClient, workerPool *workerpool.WorkerPool) *Service {
	return &Service{
		marketClient: marketClient,
		workerPool:   workerPool,
	}
}

// GetMarketPage returns a page of market items.
func (s Service) GetMarketPage(page int) ([]*marketdomain.MarketItem, error) {
	items, err := s.marketClient.GetPage(page)
	if err != nil {
		return nil, fmt.Errorf("getting the page of items: %w", err)
	}

	return items, nil
}
