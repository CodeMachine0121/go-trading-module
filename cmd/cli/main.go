package main

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	sqliterepo "transaction/internal/adapter/repository/sqlite"
	"transaction/internal/interface/cli"
	"transaction/internal/usecase/strategy"
	"transaction/pkg/logger"
)

func main() {
	// Initialize database connection
	db, err := gorm.Open(sqlite.Open("strategies.db"), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Run database migrations
	if err := sqliterepo.Migrate(db); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	// Initialize dependencies
	repo := sqliterepo.NewStrategyRepository(db)
	log := logger.NewSimpleLogger()
	svc := strategy.NewStrategyService(repo, log)

	// Create root command
	rootCmd := &cli.RootCommand{
		StrategyService: svc,
		Logger:          log,
	}

	// Execute command
	if err := rootCmd.Execute(os.Args[1:]); err != nil {
		log.Error("Command execution failed: " + err.Error())
		os.Exit(1)
	}
}
