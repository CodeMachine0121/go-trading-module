package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"transaction/internal/usecase/strategy"
	"transaction/pkg/logger"
)

var (
	createStrategyCmd *cobra.Command
	listStrategiesCmd *cobra.Command
	getStrategyCmd    *cobra.Command
	updateStrategyCmd *cobra.Command
	deleteStrategyCmd *cobra.Command
	toggleStrategyCmd *cobra.Command
)

// NewStrategyCommand creates the root strategy command with subcommands
func NewStrategyCommand(svc *strategy.StrategyService, log logger.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "strategy",
		Short: "Manage trading strategies",
		Long:  "Commands for managing price range trading strategies",
	}

	// Create command
	createStrategyCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new strategy",
		Long:  "Create a new trading strategy with buy and sell price limits",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("symbol")
			buyLower, _ := cmd.Flags().GetFloat64("buy-lower")
			sellUpper, _ := cmd.Flags().GetFloat64("sell-upper")

			if symbol == "" {
				return fmt.Errorf("symbol is required")
			}

			req := &strategy.CreateStrategyRequest{
				Symbol:    symbol,
				BuyLower:  buyLower,
				SellUpper: sellUpper,
			}

			result, err := svc.CreateStrategy(req)
			if err != nil {
				log.Error("Failed to create strategy", "error", err.Error())
				return err
			}

			log.Info("Strategy created successfully", "id", result.ID)
			fmt.Printf("Created strategy: ID=%s, Symbol=%s, BuyLower=%.2f, SellUpper=%.2f, Active=%v\n",
				result.ID, result.Symbol, result.BuyLower, result.SellUpper, result.IsActive)
			return nil
		},
	}

	createStrategyCmd.Flags().StringP("symbol", "s", "", "Symbol (e.g., BTC/USD)")
	createStrategyCmd.Flags().Float64P("buy-lower", "b", 0, "Buy lower limit")
	createStrategyCmd.Flags().Float64P("sell-upper", "u", 0, "Sell upper limit")
	createStrategyCmd.MarkFlagRequired("symbol")
	createStrategyCmd.MarkFlagRequired("buy-lower")
	createStrategyCmd.MarkFlagRequired("sell-upper")

	// List command
	listStrategiesCmd = &cobra.Command{
		Use:   "list",
		Short: "List all strategies",
		Long:  "Display all existing trading strategies",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Listing all strategies")

			results, err := svc.ListStrategies()
			if err != nil {
				log.Error("Failed to list strategies", "error", err.Error())
				return err
			}

			if len(results) == 0 {
				fmt.Println("No strategies found")
				return nil
			}

			fmt.Println("Strategies:")
			fmt.Println(strings.Repeat("-", 100))
			for _, s := range results {
				status := "Active"
				if !s.IsActive {
					status = "Inactive"
				}
				fmt.Printf("ID: %s, Symbol: %s, BuyLower: %.2f, SellUpper: %.2f, Status: %s\n",
					s.ID, s.Symbol, s.BuyLower, s.SellUpper, status)
			}
			fmt.Println(strings.Repeat("-", 100))
			return nil
		},
	}

	// Get command
	getStrategyCmd = &cobra.Command{
		Use:   "get <strategy-id>",
		Short: "Get strategy details",
		Long:  "Display details of a specific strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			log.Info("Fetching strategy", "id", id)

			result, err := svc.GetStrategy(id)
			if err != nil {
				log.Error("Failed to get strategy", "error", err.Error())
				return err
			}

			status := "Active"
			if !result.IsActive {
				status = "Inactive"
			}
			fmt.Printf("Strategy Details:\n")
			fmt.Printf("  ID: %s\n", result.ID)
			fmt.Printf("  Symbol: %s\n", result.Symbol)
			fmt.Printf("  Buy Lower: %.2f\n", result.BuyLower)
			fmt.Printf("  Sell Upper: %.2f\n", result.SellUpper)
			fmt.Printf("  Status: %s\n", status)
			return nil
		},
	}

	// Update command
	updateStrategyCmd = &cobra.Command{
		Use:   "update <strategy-id>",
		Short: "Update strategy",
		Long:  "Update an existing trading strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			buyLower, _ := cmd.Flags().GetString("buy-lower")
			sellUpper, _ := cmd.Flags().GetString("sell-upper")

			if buyLower == "" && sellUpper == "" {
				return fmt.Errorf("at least one of --buy-lower or --sell-upper is required")
			}

			// Fetch current strategy to get symbol
			current, err := svc.GetStrategy(id)
			if err != nil {
				log.Error("Strategy not found", "id", id)
				return err
			}

			buyLowerVal := current.BuyLower
			sellUpperVal := current.SellUpper

			if buyLower != "" {
				val, err := strconv.ParseFloat(buyLower, 64)
				if err != nil {
					return fmt.Errorf("invalid buy-lower value: %v", err)
				}
				buyLowerVal = val
			}

			if sellUpper != "" {
				val, err := strconv.ParseFloat(sellUpper, 64)
				if err != nil {
					return fmt.Errorf("invalid sell-upper value: %v", err)
				}
				sellUpperVal = val
			}

			req := &strategy.UpdateStrategyRequest{
				ID:        id,
				Symbol:    current.Symbol,
				BuyLower:  buyLowerVal,
				SellUpper: sellUpperVal,
			}

			result, err := svc.UpdateStrategy(req)
			if err != nil {
				log.Error("Failed to update strategy", "error", err.Error())
				return err
			}

			log.Info("Strategy updated successfully", "id", id)
			fmt.Printf("Updated strategy: ID=%s, BuyLower=%.2f, SellUpper=%.2f\n",
				result.ID, result.BuyLower, result.SellUpper)
			return nil
		},
	}

	updateStrategyCmd.Flags().StringP("buy-lower", "b", "", "Buy lower limit")
	updateStrategyCmd.Flags().StringP("sell-upper", "u", "", "Sell upper limit")

	// Delete command
	deleteStrategyCmd = &cobra.Command{
		Use:   "delete <strategy-id>",
		Short: "Delete strategy",
		Long:  "Delete an existing trading strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			log.Info("Deleting strategy", "id", id)

			err := svc.DeleteStrategy(id)
			if err != nil {
				log.Error("Failed to delete strategy", "error", err.Error())
				return err
			}

			log.Info("Strategy deleted successfully", "id", id)
			fmt.Printf("Strategy %s deleted\n", id)
			return nil
		},
	}

	// Toggle command
	toggleStrategyCmd = &cobra.Command{
		Use:   "toggle <strategy-id>",
		Short: "Toggle strategy status",
		Long:  "Enable or disable a trading strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			log.Info("Toggling strategy status", "id", id)

			result, err := svc.ToggleStrategy(id)
			if err != nil {
				log.Error("Failed to toggle strategy", "error", err.Error())
				return err
			}

			status := "Active"
			if !result.IsActive {
				status = "Inactive"
			}
			log.Info("Strategy status toggled", "id", id, "status", status)
			fmt.Printf("Strategy %s is now %s\n", id, status)
			return nil
		},
	}

	// Add subcommands to root command
	rootCmd.AddCommand(
		createStrategyCmd,
		listStrategiesCmd,
		getStrategyCmd,
		updateStrategyCmd,
		deleteStrategyCmd,
		toggleStrategyCmd,
	)

	return rootCmd
}
