package strategy

import (
	"transaction/internal/adapter/repository"
	"transaction/internal/domain"
	"transaction/pkg/logger"

	"github.com/google/uuid"
)

// StrategyService implements business logic for strategy management.
type StrategyService struct {
	repo   repository.IStrategyRepository
	logger logger.Logger
}

// NewStrategyService creates a new instance of StrategyService.
func NewStrategyService(repo repository.IStrategyRepository, logger logger.Logger) *StrategyService {
	return &StrategyService{
		repo:   repo,
		logger: logger,
	}
}

// CreateStrategy creates a new strategy.
func (s *StrategyService) CreateStrategy(req *CreateStrategyRequest) (*StrategyResponse, error) {
	s.logger.Info("Creating strategy", "symbol", req.Symbol)

	strategy := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    req.Symbol,
		BuyLower:  req.BuyLower,
		SellUpper: req.SellUpper,
		IsActive:  true,
	}

	if err := strategy.Validate(); err != nil {
		s.logger.Error("Strategy validation failed", "error", err.Error())
		return nil, err
	}

	created, err := s.repo.Create(strategy)
	if err != nil {
		s.logger.Error("Failed to create strategy", "error", err.Error())
		return nil, err
	}

	return toResponse(created), nil
}

// GetStrategy retrieves a strategy by ID.
func (s *StrategyService) GetStrategy(id string) (*StrategyResponse, error) {
	s.logger.Info("Fetching strategy", "id", id)

	strategy, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Strategy not found", "id", id)
		return nil, err
	}

	return toResponse(strategy), nil
}

// ListStrategies retrieves all strategies.
func (s *StrategyService) ListStrategies() ([]*StrategyResponse, error) {
	s.logger.Info("Listing all strategies")

	strategies, err := s.repo.FindAll()
	if err != nil {
		s.logger.Error("Failed to list strategies", "error", err.Error())
		return nil, err
	}

	responses := make([]*StrategyResponse, len(strategies))
	for i, strategy := range strategies {
		responses[i] = toResponse(strategy)
	}
	return responses, nil
}

// UpdateStrategy updates an existing strategy.
func (s *StrategyService) UpdateStrategy(req *UpdateStrategyRequest) (*StrategyResponse, error) {
	s.logger.Info("Updating strategy", "id", req.ID)

	strategy := &domain.Strategy{
		ID:        req.ID,
		Symbol:    req.Symbol,
		BuyLower:  req.BuyLower,
		SellUpper: req.SellUpper,
	}

	if err := strategy.Validate(); err != nil {
		s.logger.Error("Strategy validation failed", "error", err.Error())
		return nil, err
	}

	updated, err := s.repo.Update(strategy)
	if err != nil {
		s.logger.Error("Failed to update strategy", "id", req.ID)
		return nil, err
	}

	return toResponse(updated), nil
}

// DeleteStrategy deletes a strategy by ID.
func (s *StrategyService) DeleteStrategy(id string) error {
	s.logger.Info("Deleting strategy", "id", id)

	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete strategy", "id", id)
		return err
	}

	return nil
}

// ToggleStrategy toggles the active status of a strategy.
func (s *StrategyService) ToggleStrategy(id string) (*StrategyResponse, error) {
	s.logger.Info("Toggling strategy status", "id", id)

	strategy, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Strategy not found", "id", id)
		return nil, err
	}

	strategy.IsActive = !strategy.IsActive

	updated, err := s.repo.Update(strategy)
	if err != nil {
		s.logger.Error("Failed to toggle strategy", "id", id)
		return nil, err
	}

	return toResponse(updated), nil
}

// toResponse converts a domain Strategy to a StrategyResponse.
func toResponse(s *domain.Strategy) *StrategyResponse {
	return &StrategyResponse{
		ID:        s.ID,
		Symbol:    s.Symbol,
		BuyLower:  s.BuyLower,
		SellUpper: s.SellUpper,
		IsActive:  s.IsActive,
	}
}
