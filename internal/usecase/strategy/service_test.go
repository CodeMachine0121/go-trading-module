package strategy

import (
	"testing"
	"transaction/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of IStrategyRepository.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(strategy *domain.Strategy) (*domain.Strategy, error) {
	args := m.Called(strategy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Strategy), args.Error(1)
}

func (m *MockRepository) FindByID(id string) (*domain.Strategy, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Strategy), args.Error(1)
}

func (m *MockRepository) FindAll() ([]*domain.Strategy, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Strategy), args.Error(1)
}

func (m *MockRepository) Update(strategy *domain.Strategy) (*domain.Strategy, error) {
	args := m.Called(strategy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Strategy), args.Error(1)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockLogger is a mock implementation of Logger.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestCreateStrategy_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &CreateStrategyRequest{
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockRepo.On("Create", mock.MatchedBy(func(s *domain.Strategy) bool {
		return s.Symbol == "BTC" && s.BuyLower == 30000.0 && s.SellUpper == 50000.0 && s.IsActive
	})).Return(&domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
		IsActive:  true,
	}, nil)

	resp, err := service.CreateStrategy(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "BTC", resp.Symbol)
}

func TestCreateStrategy_InvalidPrice(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &CreateStrategyRequest{
		Symbol:    "BTC",
		BuyLower:  -100.0,
		SellUpper: 50000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	resp, err := service.CreateStrategy(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateStrategy_InvalidBoundaryRelation(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &CreateStrategyRequest{
		Symbol:    "BTC",
		BuyLower:  50000.0,
		SellUpper: 30000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()

	resp, err := service.CreateStrategy(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetStrategy_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindByID", "test-id").Return(&domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
		IsActive:  true,
	}, nil)

	resp, err := service.GetStrategy("test-id")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-id", resp.ID)
}

func TestGetStrategy_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindByID", "nonexistent").Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.GetStrategy("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestListStrategies_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	strategies := []*domain.Strategy{
		{
			ID:        "id-1",
			Symbol:    "BTC",
			BuyLower:  30000.0,
			SellUpper: 50000.0,
			IsActive:  true,
		},
		{
			ID:        "id-2",
			Symbol:    "ETH",
			BuyLower:  2000.0,
			SellUpper: 3000.0,
			IsActive:  true,
		},
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindAll").Return(strategies, nil)

	resp, err := service.ListStrategies()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 2)
}

func TestUpdateStrategy_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &UpdateStrategyRequest{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  25000.0,
		SellUpper: 55000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("Update", mock.MatchedBy(func(s *domain.Strategy) bool {
		return s.ID == "test-id" && s.BuyLower == 25000.0
	})).Return(&domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  25000.0,
		SellUpper: 55000.0,
		IsActive:  true,
	}, nil)

	resp, err := service.UpdateStrategy(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 25000.0, resp.BuyLower)
}

func TestDeleteStrategy_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("Delete", "test-id").Return(nil)

	err := service.DeleteStrategy("test-id")
	assert.NoError(t, err)
}

func TestToggleStrategy_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	strategy := &domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
		IsActive:  true,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindByID", "test-id").Return(strategy, nil)
	mockRepo.On("Update", mock.MatchedBy(func(s *domain.Strategy) bool {
		return s.ID == "test-id" && !s.IsActive
	})).Return(&domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
		IsActive:  false,
	}, nil)

	resp, err := service.ToggleStrategy("test-id")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsActive)
}

func TestCreateStrategy_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &CreateStrategyRequest{
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("Create", mock.MatchedBy(func(s *domain.Strategy) bool {
		return s.Symbol == "BTC"
	})).Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.CreateStrategy(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateStrategy_StrategyNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	req := &UpdateStrategyRequest{
		ID:        "non-existent",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("Update", mock.MatchedBy(func(s *domain.Strategy) bool {
		return s.ID == "non-existent"
	})).Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.UpdateStrategy(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteStrategy_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("Delete", "test-id").Return(domain.ErrStrategyNotFound)

	err := service.DeleteStrategy("test-id")
	assert.Error(t, err)
}

func TestListStrategies_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindAll").Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.ListStrategies()
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestToggleStrategy_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindByID", "non-existent").Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.ToggleStrategy("non-existent")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestToggleStrategy_UpdateError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := new(MockLogger)
	service := NewStrategyService(mockRepo, mockLogger)

	strategy := &domain.Strategy{
		ID:        "test-id",
		Symbol:    "BTC",
		BuyLower:  30000.0,
		SellUpper: 50000.0,
		IsActive:  true,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockRepo.On("FindByID", "test-id").Return(strategy, nil)
	mockRepo.On("Update", mock.Anything).Return(nil, domain.ErrStrategyNotFound)

	resp, err := service.ToggleStrategy("test-id")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
