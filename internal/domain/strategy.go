package domain

import (
	"errors"
	"time"
)

// Strategy represents a price range strategy for cryptocurrency trading.
type Strategy struct {
	ID        string  `gorm:"primaryKey"`
	Symbol    string  // BTC, ETH, USDT, etc.
	BuyLower  float64 // Minimum price to trigger buy signal
	SellUpper float64 // Maximum price to trigger sell signal
	IsActive  bool    // Whether the strategy is currently active
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate checks if the strategy has valid configuration.
func (s *Strategy) Validate() error {
	if s.BuyLower <= 0 {
		return errors.New("buy lower bound must be positive")
	}
	if s.SellUpper <= s.BuyLower {
		return errors.New("sell upper bound must be greater than buy lower bound")
	}
	return nil
}

// ShouldBuy determines if the current price triggers a buy signal.
func (s *Strategy) ShouldBuy(currentPrice float64) bool {
	return s.IsActive && currentPrice <= s.BuyLower
}

// ShouldSell determines if the current price triggers a sell signal.
func (s *Strategy) ShouldSell(currentPrice float64) bool {
	return s.IsActive && currentPrice >= s.SellUpper
}
