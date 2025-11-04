package strategy

// CreateStrategyRequest represents the request to create a new strategy.
type CreateStrategyRequest struct {
	Symbol    string  // BTC, ETH, USDT, etc.
	BuyLower  float64 // Minimum price to trigger buy signal
	SellUpper float64 // Maximum price to trigger sell signal
}

// UpdateStrategyRequest represents the request to update an existing strategy.
type UpdateStrategyRequest struct {
	ID        string  // Strategy ID
	Symbol    string  // BTC, ETH, USDT, etc.
	BuyLower  float64 // Minimum price to trigger buy signal
	SellUpper float64 // Maximum price to trigger sell signal
}

// StrategyResponse represents the response containing strategy data.
type StrategyResponse struct {
	ID        string  // Unique identifier
	Symbol    string  // BTC, ETH, USDT, etc.
	BuyLower  float64 // Minimum price to trigger buy signal
	SellUpper float64 // Maximum price to trigger sell signal
	IsActive  bool    // Whether the strategy is currently active
}
