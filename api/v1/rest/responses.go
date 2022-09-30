package rest

import "github.com/UArt-project/UArt-proxy/domain/marketdomain"

type MarketPageResponse struct {
	// The page of the market items.
	Page int `json:"page"`
	// Items.
	Items []*MarketItemResponse `json:"items"`
}

type MarketItemResponse struct {
	// The ID of the item.
	ID string `json:"id"`
	// The name of the item.
	Name string `json:"name"`
	// The price of the item.
	Price float64 `json:"price"`
	// Photo of the item.
	Photo string `json:"photo"`
}

// itemsToResponse converts the market items to the response.
func itemsToResponse(page int, items []marketdomain.MarketItem) *MarketPageResponse {
	var returnItems []*MarketItemResponse

	for _, item := range items {
		returnItems = append(returnItems, &MarketItemResponse{
			ID:    item.ID,
			Name:  item.Name,
			Price: item.Price,
			Photo: item.Photo,
		})
	}

	result := &MarketPageResponse{
		Page:  page,
		Items: returnItems,
	}

	return result
}
