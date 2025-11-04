package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		buyLower  float64
		sellUpper float64
		wantErr   bool
	}{
		{
			name:      "should create valid strategy",
			symbol:    "BTC",
			buyLower:  60000,
			sellUpper: 70000,
			wantErr:   false,
		},
		{
			name:      "should fail when buy lower is zero",
			symbol:    "BTC",
			buyLower:  0,
			sellUpper: 70000,
			wantErr:   true,
		},
		{
			name:      "should fail when buy lower is negative",
			symbol:    "BTC",
			buyLower:  -60000,
			sellUpper: 70000,
			wantErr:   true,
		},
		{
			name:      "should fail when sell upper <= buy lower",
			symbol:    "BTC",
			buyLower:  70000,
			sellUpper: 70000,
			wantErr:   true,
		},
		{
			name:      "should fail when sell upper < buy lower",
			symbol:    "BTC",
			buyLower:  70000,
			sellUpper: 60000,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := &Strategy{
				ID:        "test-id",
				Symbol:    tt.symbol,
				BuyLower:  tt.buyLower,
				SellUpper: tt.sellUpper,
				IsActive:  true,
				CreatedAt: time.Now(),
			}

			err := strategy.Validate()

			if tt.wantErr {
				assert.Error(t, err, "expected validation error")
			} else {
				assert.NoError(t, err, "expected no validation error")
			}
		})
	}
}

func TestStrategyValidate(t *testing.T) {
	tests := []struct {
		name     string
		strategy *Strategy
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid strategy with positive prices",
			strategy: &Strategy{
				ID:        "test-1",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid strategy - buy lower must be positive",
			strategy: &Strategy{
				ID:        "test-2",
				Symbol:    "BTC",
				BuyLower:  0,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid strategy - sell upper must be > buy lower",
			strategy: &Strategy{
				ID:        "test-3",
				Symbol:    "BTC",
				BuyLower:  70000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.strategy.Validate()

			if tt.wantErr {
				assert.Error(t, err, "expected validation error")
			} else {
				assert.NoError(t, err, "expected no validation error")
			}
		})
	}
}

func TestStrategyShouldBuy(t *testing.T) {
	tests := []struct {
		name         string
		strategy     *Strategy
		currentPrice float64
		expected     bool
	}{
		{
			name: "should buy when price <= buy lower and strategy is active",
			strategy: &Strategy{
				ID:        "test-1",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 59000,
			expected:     true,
		},
		{
			name: "should buy when price exactly equals buy lower and strategy is active",
			strategy: &Strategy{
				ID:        "test-2",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 60000,
			expected:     true,
		},
		{
			name: "should not buy when price > buy lower",
			strategy: &Strategy{
				ID:        "test-3",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 61000,
			expected:     false,
		},
		{
			name: "should not buy when strategy is inactive",
			strategy: &Strategy{
				ID:        "test-4",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  false,
				CreatedAt: time.Now(),
			},
			currentPrice: 59000,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.strategy.ShouldBuy(tt.currentPrice)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStrategyShouldSell(t *testing.T) {
	tests := []struct {
		name         string
		strategy     *Strategy
		currentPrice float64
		expected     bool
	}{
		{
			name: "should sell when price >= sell upper and strategy is active",
			strategy: &Strategy{
				ID:        "test-1",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 71000,
			expected:     true,
		},
		{
			name: "should sell when price exactly equals sell upper and strategy is active",
			strategy: &Strategy{
				ID:        "test-2",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 70000,
			expected:     true,
		},
		{
			name: "should not sell when price < sell upper",
			strategy: &Strategy{
				ID:        "test-3",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  true,
				CreatedAt: time.Now(),
			},
			currentPrice: 69000,
			expected:     false,
		},
		{
			name: "should not sell when strategy is inactive",
			strategy: &Strategy{
				ID:        "test-4",
				Symbol:    "BTC",
				BuyLower:  60000,
				SellUpper: 70000,
				IsActive:  false,
				CreatedAt: time.Now(),
			},
			currentPrice: 71000,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.strategy.ShouldSell(tt.currentPrice)
			assert.Equal(t, tt.expected, result)
		})
	}
}
