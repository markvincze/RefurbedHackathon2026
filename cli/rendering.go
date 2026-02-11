package main

import (
	"fmt"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
	"github.com/yarlson/tap"
)

type stat struct {
	date    string
	revenue float64
}

var blockStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")).
	Background(lipgloss.Color("3"))

func renderRevenueByDay() {
	fmt.Println("Revenue by day")
	fmt.Println()
	data := []stat{
		stat{"2026.01.01", 120},
		stat{"2026.01.02", 112},
		stat{"2026.01.03", 270},
		stat{"2026.01.04", 178},
		stat{"2026.01.05", 153},
		stat{"2026.01.06", 240},
		stat{"2026.01.07", 225},
	}

	renderRevenueByDayTable(data)
	renderRevenueByDayGraph(data)
}

func renderRevenueByDayTable(data []stat) {
	textData := make([][]string, 0)
	for _, d := range data {
		textData = append(textData, []string{d.date, fmt.Sprintf("€ %f", d.revenue)})
	}

	tap.Table(
		[]string{"Day", "Revenue"},
		textData,
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderRevenueByDayGraph(data []stat) {
	values := make([]barchart.BarData, 0)
	for _, dayAndValue := range data {
		values = append(
			values,
			barchart.BarData{
				Label:  dayAndValue.date,
				Values: []barchart.BarValue{{Name: "Revenue", Value: dayAndValue.revenue, Style: blockStyle}}})
	}

	bc := barchart.New(140, 15)
	bc.SetShowAxis(true)
	bc.PushAll(values)
	bc.Draw()

	fmt.Println(bc.View())
}
