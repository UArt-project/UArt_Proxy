// Package marketclient provides a client for the market service.
package marketclient

import (
	"fmt"
	"io"
	"net/http"

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
}

// NewMarketServiceClient creates a new instance of the MarketServiceClient.
func NewMarketServiceClient(url string) *MarketServiceClient {
	return &MarketServiceClient{
		url: url,
	}
}

// GetPage returns a page of market items.
func (c MarketServiceClient) GetPage(page int) ([]marketdomain.MarketItem, error) {
	httpClient := http.Client{}

	req, err := http.NewRequest("GET", c.url+"/marketplace/v1/items/0", nil)
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

	err = jsonoperations.Decode(body, &items)
	if err != nil {
		return nil, fmt.Errorf("decoding the response body: %w", err)
	}

	items = append(items, marketdomain.MarketItem{ID: "-1", Name: "На Бандерасмузі", Price: 100, Photo: ""})
	items = append(items, marketdomain.MarketItem{ID: "-2", Name: "На велику бавовну", Price: 500, Photo: ""})
	items = append(items, marketdomain.MarketItem{ID: "-3", Name: "На зірку смерті", Price: 1000, Photo: ""})

	return items, nil
}
