package sqlite

import (
	"gorm.io/gorm"
	"transaction/internal/adapter/repository"
	"transaction/internal/domain"
)

// StrategyRepository implements the IStrategyRepository interface using SQLite via GORM.
type StrategyRepository struct {
	db *gorm.DB
}

// NewStrategyRepository creates a new SQLite-backed IStrategyRepository.
func NewStrategyRepository(db *gorm.DB) repository.IStrategyRepository {
	return &StrategyRepository{db: db}
}

// Create persists a new strategy and returns the created strategy.
func (r *StrategyRepository) Create(strategy *domain.Strategy) (*domain.Strategy, error) {
	result := r.db.Create(strategy)
	if result.Error != nil {
		return nil, result.Error
	}
	return strategy, nil
}

// FindByID retrieves a strategy by its ID.
func (r *StrategyRepository) FindByID(id string) (*domain.Strategy, error) {
	strategy := &domain.Strategy{}
	result := r.db.First(strategy, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, domain.ErrStrategyNotFound
		}
		return nil, result.Error
	}
	return strategy, nil
}

// FindAll retrieves all strategies from the database.
func (r *StrategyRepository) FindAll() ([]*domain.Strategy, error) {
	strategies := make([]*domain.Strategy, 0)
	result := r.db.Find(&strategies)
	if result.Error != nil {
		return nil, result.Error
	}
	return strategies, nil
}

// Update modifies an existing strategy.
func (r *StrategyRepository) Update(strategy *domain.Strategy) (*domain.Strategy, error) {
	result := r.db.Save(strategy)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrStrategyNotFound
	}
	return strategy, nil
}

// Delete removes a strategy by its ID.
func (r *StrategyRepository) Delete(id string) error {
	result := r.db.Delete(&domain.Strategy{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrStrategyNotFound
	}
	return nil
}
