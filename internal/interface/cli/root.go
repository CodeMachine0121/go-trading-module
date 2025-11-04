package cli

import (
	"github.com/spf13/cobra"
	"transaction/internal/usecase/strategy"
	"transaction/pkg/logger"
)

// RootCommand is the root CLI command
type RootCommand struct {
	StrategyService *strategy.StrategyService
	Logger          logger.Logger
}

// Execute runs the CLI application
func (r *RootCommand) Execute(args []string) error {
	rootCmd := &cobra.Command{
		Use:   "strategy-cli",
		Short: "Trading Strategy Management CLI",
		Long:  "A command-line tool for managing cryptocurrency trading strategies",
	}

	// Add strategy command
	strategyCmd := NewStrategyCommand(r.StrategyService, r.Logger)
	rootCmd.AddCommand(strategyCmd)

	// Set args
	rootCmd.SetArgs(args)

	// Execute
	return rootCmd.Execute()
}
