package repository

import "transaction/internal/domain"

// IStrategyRepository defines the interface for persisting Strategy entities.
type IStrategyRepository interface {
	// Create persists a new strategy and returns the created strategy with ID.
	Create(strategy *domain.Strategy) (*domain.Strategy, error)

	// FindByID retrieves a strategy by its ID.
	// Returns ErrStrategyNotFound if the strategy does not exist.
	FindByID(id string) (*domain.Strategy, error)

	// FindAll retrieves all strategies from the repository.
	FindAll() ([]*domain.Strategy, error)

	// Update modifies an existing strategy.
	// Returns ErrStrategyNotFound if the strategy does not exist.
	Update(strategy *domain.Strategy) (*domain.Strategy, error)

	// Delete removes a strategy by its ID.
	// Returns ErrStrategyNotFound if the strategy does not exist.
	Delete(id string) error
}
