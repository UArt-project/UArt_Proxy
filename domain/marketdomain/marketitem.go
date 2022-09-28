// Package marketdomain provides objects used with the market service.
package marketdomain

// MarketItem represents a market item.
type MarketItem struct {
	// The ID of the item.
	ID string `json:"id"`
	// The name of the item.
	Name string `json:"name"`
	// The price of the item.
	Price float64 `json:"price"`
	// Photo of the item.
	Photo string `json:"photoLink"`
}
