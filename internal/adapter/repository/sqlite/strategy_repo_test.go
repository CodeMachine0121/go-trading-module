package sqlite

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"transaction/internal/domain"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "failed to open in-memory database")

	// Run migrations
	err = Migrate(db)
	require.NoError(t, err, "failed to run migrations")

	return db
}

func TestCreate_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategy := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "BTC",
		BuyLower:  40000.0,
		SellUpper: 60000.0,
		IsActive:  true,
	}

	created, err := repo.Create(strategy)
	require.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, strategy.ID, created.ID)
	assert.Equal(t, "BTC", created.Symbol)
	assert.Equal(t, 40000.0, created.BuyLower)
	assert.Equal(t, 60000.0, created.SellUpper)
	assert.True(t, created.IsActive)
}

func TestFindByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategy := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "ETH",
		BuyLower:  2000.0,
		SellUpper: 4000.0,
		IsActive:  true,
	}
	_, err := repo.Create(strategy)
	require.NoError(t, err)

	found, err := repo.FindByID(strategy.ID)
	require.NoError(t, err)
	assert.Equal(t, strategy.ID, found.ID)
	assert.Equal(t, "ETH", found.Symbol)
}

func TestFindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	_, err := repo.FindByID("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrStrategyNotFound, err)
}

func TestFindAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategies, err := repo.FindAll()
	require.NoError(t, err)
	assert.Empty(t, strategies)
}

func TestFindAll_Multiple(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategy1 := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "BTC",
		BuyLower:  40000.0,
		SellUpper: 60000.0,
		IsActive:  true,
	}
	strategy2 := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "ETH",
		BuyLower:  2000.0,
		SellUpper: 4000.0,
		IsActive:  false,
	}

	_, err := repo.Create(strategy1)
	require.NoError(t, err)
	_, err = repo.Create(strategy2)
	require.NoError(t, err)

	strategies, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, strategies, 2)
}

func TestUpdate_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategy := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "BTC",
		BuyLower:  40000.0,
		SellUpper: 60000.0,
		IsActive:  true,
	}
	_, err := repo.Create(strategy)
	require.NoError(t, err)

	// Update the strategy
	strategy.BuyLower = 35000.0
	strategy.SellUpper = 65000.0
	strategy.IsActive = false

	updated, err := repo.Update(strategy)
	require.NoError(t, err)
	assert.Equal(t, 35000.0, updated.BuyLower)
	assert.Equal(t, 65000.0, updated.SellUpper)
	assert.False(t, updated.IsActive)
}

func TestDelete_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewStrategyRepository(db)

	strategy := &domain.Strategy{
		ID:        uuid.New().String(),
		Symbol:    "BTC",
		BuyLower:  40000.0,
		SellUpper: 60000.0,
		IsActive:  true,
	}
	_, err := repo.Create(strategy)
	require.NoError(t, err)

	// Verify it exists
	_, err = repo.FindByID(strategy.ID)
	require.NoError(t, err)

	// Delete it
	err = repo.Delete(strategy.ID)
	require.NoError(t, err)

	// Verify it's gone
	_, err = repo.FindByID(strategy.ID)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrStrategyNotFound, err)
}
