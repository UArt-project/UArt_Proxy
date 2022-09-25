// Package marketdomain provides objects used with the market service.
package marketdomain

// MarketItem represents a market item.
type MarketItem struct {
	// The ID of the item.
	ID string
	// The name of the item.
	Name string
	// The price of the item.
	Price float64
	// Photo of the item.
	Photo string
}
