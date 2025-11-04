package sqlite

import (
	"gorm.io/gorm"
	"transaction/internal/domain"
)

// Migrate runs all database migrations.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Strategy{})
}

// RunMigration is an alias for Migrate for convenience.
func RunMigration(db *gorm.DB) error {
	return Migrate(db)
}
