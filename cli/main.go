package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yarlson/tap"
	"refurbed.com/hackathon/reporting"
)

func main() {
	var dataset *reporting.OrderDataset

	for {
		clearScreen()
		tap.Intro("Welcome to order data visualizer!")

		if dataset == nil {
			spinner := tap.NewSpinner(tap.SpinnerOptions{})
			spinner.Start("Loading the dataset...")

			in, err := os.Open("orders_v3.csv")

			if err != nil {
				fmt.Printf("Opening the file failed: %v", err)
				return
			}

			dataset, err = reporting.ImportOrderDatasetFromCSV(in)

			if err != nil {
				fmt.Printf("Processing the data failed: %v", err)
				return
			}

			spinner.Stop("Loading complete", 0)
		}

		tap.Message(fmt.Sprintf("AOV: %v", dataset.AOV()))
		tap.Message("Total revenue: TODO")
		tap.Message("Delivery days median: TODO, p95: TODO")
		tap.Message("Return rate: TODO")

		options := []tap.SelectOption[string]{
			{Value: "RevenueByDay", Label: "Revenue by day", Hint: ""},
			{Value: "RevenueByWeek", Label: "Revenue by week", Hint: ""},
			{Value: "ReturnRateByCategory", Label: "Return rate by category", Hint: ""},
			{Value: "OrderCountByCategory", Label: "Order count by category and subcategory", Hint: ""},
			{Value: "QueryBuilder", Label: "Query builder", Hint: "Create custom query"},
			{Value: "Quit", Label: "Quit"},
		}

		result := tap.Select[string](context.Background(), tap.SelectOptions[string]{
			Message: "Select from the below options:",
			Options: options,
		})

		clearScreen()

		switch result {
		case "RevenueByDay":
			renderRevenueByDay()
		case "RevenueByWeek":
			renderRevenueByWeek()
		case "ReturnRateByCategory":
			renderReturnRateByCategory(dataset)
		case "Quit":
			return
		}

		tap.Text(context.Background(), tap.TextOptions{Message: "Press Enter to return to menu"})
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
