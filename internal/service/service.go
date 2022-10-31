// Package service contains main logic of the application.
package service

import (
	"fmt"

	"github.com/UArt-project/UArt-proxy/domain/authdomain"
	"github.com/UArt-project/UArt-proxy/domain/marketdomain"
	"github.com/UArt-project/UArt-proxy/pkg/cache"
	"github.com/UArt-project/UArt-proxy/pkg/clients/authclient"
	"github.com/UArt-project/UArt-proxy/pkg/clients/marketclient"
	"github.com/UArt-project/UArt-proxy/pkg/workerpool"
)

// AppService provides information of main application service functionality.
type AppService interface {
	// GetMarketPage returns a page of market items.
	GetMarketPage(page int) ([]marketdomain.MarketItem, error)

	GetAuthPage() (string, error)

	GetAuthToken(callbackData authdomain.CallbackRequest) (*authdomain.AuthReturn, error)
}

// Service is a main application logic.
type Service struct {
	// The market client.
	marketClient marketclient.MarketClient
	// The auth client.
	authClient authclient.AuthClient
	// Worker pool.
	workerPool *workerpool.WorkerPool
	// The cache.
	cache *cache.LocalCache
}

// NewService creates a new instance of the Service.
func NewService(marketClient marketclient.MarketClient, authClient authclient.AuthClient,
	workerPool *workerpool.WorkerPool, cache *cache.LocalCache,
) *Service {
	return &Service{
		marketClient: marketClient,
		authClient:   authClient,
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

// GetAuthPage returns a auth redirection page.
func (s Service) GetAuthPage() (string, error) {
	url, err := s.authClient.SendAuthRequest()
	if err != nil {
		return "", fmt.Errorf("getting the auth page: %w", err)
	}

	return url, nil
}

// GetAuthToken returns the OAuth data.
func (s Service) GetAuthToken(callbackData authdomain.CallbackRequest) (*authdomain.AuthReturn, error) {
	data, err := s.authClient.SendOAuthData(callbackData)
	if err != nil {
		return nil, fmt.Errorf("getting the OAuth data: %w", err)
	}

	return data, nil
}
