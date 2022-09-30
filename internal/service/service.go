// Package service contains main logic of the application.
package service

import (
	"fmt"

	"github.com/UArt-project/UArt-proxy/domain/marketdomain"
	"github.com/UArt-project/UArt-proxy/pkg/cache"
	"github.com/UArt-project/UArt-proxy/pkg/clients/marketclient"
	"github.com/UArt-project/UArt-proxy/pkg/workerpool"
)

// AppService provides information of main application service functionality.
type AppService interface {
	// GetMarketPage returns a page of market items.
	GetMarketPage(page int) ([]marketdomain.MarketItem, error)
}

// Service is a main application logic.
type Service struct {
	// The market client.
	marketClient marketclient.MarketClient
	// Worker pool.
	workerPool *workerpool.WorkerPool
	// The cache.
	cache *cache.LocalCache
}

// NewService creates a new instance of the Service.
func NewService(marketClient marketclient.MarketClient, workerPool *workerpool.WorkerPool, cache *cache.LocalCache) *Service {
	return &Service{
		marketClient: marketClient,
		workerPool:   workerPool,
		cache:        cache,
	}
}

// GetMarketPage returns a page of market items.
func (s Service) GetMarketPage(page int) ([]marketdomain.MarketItem, error) {
	items, err := s.cache.Read(page)
	if err == nil && len(items) > 0 {
		return items, nil
	}

	items, err = s.marketClient.GetPage(page)
	if err != nil {
		return nil, fmt.Errorf("getting the page of items: %w", err)
	}

	return items, nil
}
