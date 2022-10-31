// Package marketclient provides a client for the market service.
package marketclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/UArt-project/UArt-proxy/domain/marketdomain"
	"github.com/UArt-project/UArt-proxy/pkg/jsonoperations"
)

type MarketClient interface {
	// GetPage returns a page of market items.
	GetPage(page int) ([]marketdomain.MarketItem, error)
}

// MarketServiceClient is a client for the market service.
type MarketServiceClient struct {
	// The url of the market service.
	url string
	// Timeout for the request.
	timeout time.Duration
}

// NewMarketServiceClient creates a new instance of the MarketServiceClient.
func NewMarketServiceClient(url string, timeout time.Duration) *MarketServiceClient {
	return &MarketServiceClient{
		url:     url,
		timeout: timeout,
	}
}

// GetPage returns a page of market items.
func (c MarketServiceClient) GetPage(page int) ([]marketdomain.MarketItem, error) {
	var httpClient http.Client

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/marketplace/v1/items/"+strconv.Itoa(page), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for getting the page of items: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getting the page of items: %w", err)
	}

	var items []marketdomain.MarketItem

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading the response body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing the response body: %w", err)
	}

	err = jsonoperations.Decode(body, &items)
	if err != nil {
		return nil, fmt.Errorf("decoding the response body: %w", err)
	}

	banderaSmoothiePrice := 100.0
	bawownaPrice := 500.0
	deathStarPrice := 1000.0

	items = append(items, marketdomain.MarketItem{
		ID:    "-1",
		Name:  "На Бандерасмузі",
		Price: banderaSmoothiePrice,
		Photo: "",
	})
	items = append(items, marketdomain.MarketItem{
		ID:    "-2",
		Name:  "На велику бавовну",
		Price: bawownaPrice,
		Photo: "",
	})
	items = append(items, marketdomain.MarketItem{
		ID:    "-3",
		Name:  "На зірку смерті",
		Price: deathStarPrice,
		Photo: "",
	})

	return items, nil
}
