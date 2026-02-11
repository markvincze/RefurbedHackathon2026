package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/yarlson/tap"
)

type model struct {
	chart       tslc.Model
	zoneManager *zone.Manager
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// forward Bubble Tea Msg to time series chart
	// and draw all data sets using braille runes
	m.chart, _ = m.chart.Update(msg)
	m.chart.DrawBrailleAll()
	return m, nil
}

func (m model) View() string {
	// call bubblezone Manager.Scan() at root model
	return m.zoneManager.Scan(
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")). // purple
			Render(m.chart.View()),
	)
}

func main() {
	for {
		clearScreen()
		tap.Intro("Welcome to order data visualizer!")

		spinner := tap.NewSpinner(tap.SpinnerOptions{})
		spinner.Start("Loading the dataset...")
		time.Sleep(100 * time.Millisecond)
		spinner.Stop("Loading complete", 0)

		tap.Message("AOV: TODO")
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
			renderReturnRateByCategory()
		case "Quit":
			return
		}

		tap.Text(context.Background(), tap.TextOptions{Message: "Press Enter to return to menu"})
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func drawGraph() {
	width := 30
	height := 12
	chart := tslc.New(width, height)

	// add default data set
	dataSet := []float64{0, 2, 4, 6, 8, 10, 8, 6, 4, 2, 0}
	for i, v := range dataSet {
		date := time.Now().Add(time.Hour * time.Duration(24*i))
		chart.Push(tslc.TimePoint{date, v})
	}

	// set default data set line color to red
	chart.SetStyle(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")), // red
	)

	// add additional data set by name
	dataSet2 := []float64{10, 8, 6, 4, 2, 0, 2, 4, 6, 8, 10}
	for i, v := range dataSet2 {
		date := time.Now().Add(time.Hour * time.Duration(24*i))
		chart.PushDataSet("dataSet2", tslc.TimePoint{date, v})
	}

	// set additional data set line color to green
	chart.SetDataSetStyle("dataSet2",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")), // green
	)

	// mouse support is enabled with BubbleZone
	zoneManager := zone.New()
	chart.SetZoneManager(zoneManager)
	chart.Focus() // set focus to process keyboard and mouse messages

	// start new Bubble Tea program with mouse support enabled
	m := model{chart, zoneManager}
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
