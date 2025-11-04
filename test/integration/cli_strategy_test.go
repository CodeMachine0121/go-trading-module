package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	sqliterepo "transaction/internal/adapter/repository/sqlite"
	"transaction/internal/usecase/strategy"
	"transaction/pkg/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, sqliterepo.RunMigration(db))
	return db
}

// TestCLICreateStrategy tests the strategy create command
func TestCLICreateStrategy(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create a valid strategy using the service
	createReq := &strategy.CreateStrategyRequest{
		Symbol:    "BTC/USD",
		BuyLower:  50000,
		SellUpper: 60000,
	}

	result, err := svc.CreateStrategy(createReq)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTC/USD", result.Symbol)
	assert.Equal(t, float64(50000), result.BuyLower)
	assert.Equal(t, float64(60000), result.SellUpper)
	assert.True(t, result.IsActive)
}

// TestCLIListStrategies tests the strategy list command
func TestCLIListStrategies(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create multiple strategies
	for i := 0; i < 3; i++ {
		_, err := svc.CreateStrategy(&strategy.CreateStrategyRequest{
			Symbol:    "BTC/USD",
			BuyLower:  50000 + float64(i*1000),
			SellUpper: 60000 + float64(i*1000),
		})
		require.NoError(t, err)
	}

	// List strategies
	results, err := svc.ListStrategies()
	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

// TestCLIGetStrategy tests the strategy get command
func TestCLIGetStrategy(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create a strategy
	created, err := svc.CreateStrategy(&strategy.CreateStrategyRequest{
		Symbol:    "ETH/USD",
		BuyLower:  3000,
		SellUpper: 4000,
	})
	require.NoError(t, err)

	// Get the strategy
	result, err := svc.GetStrategy(created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, "ETH/USD", result.Symbol)
}

// TestCLIUpdateStrategy tests the strategy update command
func TestCLIUpdateStrategy(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create a strategy
	created, err := svc.CreateStrategy(&strategy.CreateStrategyRequest{
		Symbol:    "BTC/USD",
		BuyLower:  50000,
		SellUpper: 60000,
	})
	require.NoError(t, err)

	// Update the strategy
	updateReq := &strategy.UpdateStrategyRequest{
		ID:        created.ID,
		BuyLower:  51000,
		SellUpper: 61000,
	}
	result, err := svc.UpdateStrategy(updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, float64(51000), result.BuyLower)
	assert.Equal(t, float64(61000), result.SellUpper)
}

// TestCLIDeleteStrategy tests the strategy delete command
func TestCLIDeleteStrategy(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create a strategy
	created, err := svc.CreateStrategy(&strategy.CreateStrategyRequest{
		Symbol:    "BTC/USD",
		BuyLower:  50000,
		SellUpper: 60000,
	})
	require.NoError(t, err)

	// Delete the strategy
	err = svc.DeleteStrategy(created.ID)
	assert.NoError(t, err)

	// Verify it's deleted
	result, err := svc.GetStrategy(created.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestCLIToggleStrategy tests the strategy toggle command
func TestCLIToggleStrategy(t *testing.T) {
	db := setupTestDB(t)
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create a strategy (should be active by default)
	created, err := svc.CreateStrategy(&strategy.CreateStrategyRequest{
		Symbol:    "BTC/USD",
		BuyLower:  50000,
		SellUpper: 60000,
	})
	require.NoError(t, err)
	assert.True(t, created.IsActive)

	// Toggle to inactive
	result, err := svc.ToggleStrategy(created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsActive, "First toggle should set IsActive to false")

	// Toggle back to active
	result, err = svc.ToggleStrategy(created.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsActive, "Second toggle should set IsActive to true")
}

// TestLoggerOutput tests that the logger produces output
func TestLoggerOutput(t *testing.T) {
	log := logger.NewSimpleLogger()

	// Test that logger methods don't panic
	assert.NotPanics(t, func() {
		log.Info("Test info message with %s", "arguments")
		log.Error("Test error message with %d", 42)
		log.Warn("Test warn message")
	})
}
