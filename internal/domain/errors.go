package domain

import "errors"

// Domain-level errors for strategy management.
var (
	// ErrInvalidStrategy indicates that the strategy configuration is invalid.
	ErrInvalidStrategy = errors.New("invalid strategy")

	// ErrInvalidPrice indicates that the price value is invalid.
	ErrInvalidPrice = errors.New("invalid price")

	// ErrStrategyNotFound indicates that the requested strategy does not exist.
	ErrStrategyNotFound = errors.New("strategy not found")
)
