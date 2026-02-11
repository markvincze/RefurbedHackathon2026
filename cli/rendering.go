package main

import (
	"fmt"
	"time"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
	"github.com/yarlson/tap"
	"refurbed.com/hackathon/reporting"
)

type stat struct {
	x string
	y float64
}

var blockStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")).
	Background(lipgloss.Color("3"))

func renderRevenueByDay(dataset *reporting.OrderDataset) {
	fmt.Println("Revenue by day")
	fmt.Println()

	today := time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC)

	data := make([]stat, 0)
	rbd := dataset.RevenueByDay(today.AddDate(0, 0, -8), today.AddDate(0, 0, -1))

	for _, r := range rbd {
		data = append(data, stat{r.Title, r.Revenue.InexactFloat64()})
	}

	renderRevenueByDayTable(data)
	renderRevenueByDayGraph(data)
}

func renderRevenueByDayTable(data []stat) {
	textData := make([][]string, 0)
	for _, d := range data {
		textData = append(textData, []string{d.x, fmt.Sprintf("€ %f", d.y)})
	}

	tap.Table(
		[]string{"Day", "Revenue"},
		textData,
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderRevenueByDayGraph(data []stat) {
	values := make([]barchart.BarData, 0)
	for _, stat := range data {
		values = append(
			values,
			barchart.BarData{
				Label:  stat.x,
				Values: []barchart.BarValue{{Name: "Revenue", Value: stat.y, Style: blockStyle}}})
	}

	bc := barchart.New(140, 15)
	bc.SetShowAxis(true)
	bc.PushAll(values)
	bc.Draw()

	fmt.Println(bc.View())
}

func renderRevenueByWeek(dataset *reporting.OrderDataset) {
	fmt.Println("Revenue by week")
	fmt.Println()
	today := time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC)

	data := make([]stat, 0)
	rbd := dataset.RevenueByWeek(today.AddDate(0, 0, -(7*7)-1), today.AddDate(0, 0, -1))

	for _, r := range rbd {
		data = append(data, stat{r.Title, r.Revenue.InexactFloat64()})
	}

	// data := []stat{
	// 	{"2026.02.09", 120},
	// 	{"2026.02.02", 112},
	// 	{"2026.01.26", 270},
	// 	{"2026.01.19", 178},
	// 	{"2026.01.12", 153},
	// 	{"2026.01.5", 240},
	// 	{"2025.12.29", 225},
	// }

	renderRevenueByWeekTable(data)
	renderRevenueByWeekGraph(data)
}

func renderRevenueByWeekTable(data []stat) {
	textData := make([][]string, 0)
	for _, d := range data {
		textData = append(textData, []string{d.x, fmt.Sprintf("€ %f", d.y)})
	}

	tap.Table(
		[]string{"Week", "Revenue"},
		textData,
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderRevenueByWeekGraph(data []stat) {
	values := make([]barchart.BarData, 0)
	for _, stat := range data {
		values = append(
			values,
			barchart.BarData{
				Label:  stat.x,
				Values: []barchart.BarValue{{Name: "Revenue", Value: stat.y, Style: blockStyle}}})
	}

	bc := barchart.New(140, 15)
	bc.SetShowAxis(true)
	bc.PushAll(values)
	bc.Draw()

	fmt.Println(bc.View())
}

func renderReturnRateByCategory(dataset *reporting.OrderDataset) {
	fmt.Println("Return rate by category")
	fmt.Println()

	data := make([]stat, 0)

	for _, c := range dataset.AllCategories() {
		data = append(data, stat{string(c), dataset.ReturnRateByCategory(c)})
	}

	renderReturnRateByCategoryTable(data)
	renderReturnRateByCategoryGraph(data)
}

func renderReturnRateByCategoryTable(data []stat) {
	textData := make([][]string, 0)
	for _, d := range data {
		textData = append(textData, []string{d.x, fmt.Sprintf("%d%%", (int)(100*d.y))})
	}

	tap.Table(
		[]string{"Category", "Return rate"},
		textData,
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderReturnRateByCategoryGraph(data []stat) {
	values := make([]barchart.BarData, 0)
	for _, stat := range data {
		values = append(
			values,
			barchart.BarData{
				Label:  stat.x,
				Values: []barchart.BarValue{{Name: "Return rate", Value: stat.y, Style: blockStyle}}})
	}

	bc := barchart.New(
		140, 15,
		barchart.WithDataSet(values))

	// bc := barchart.New(
	// 	140, len(data)*2,
	// 	barchart.WithDataSet(values),
	// 	barchart.WithHorizontalBars())

	bc.Draw()

	fmt.Println(bc.View())
}

func renderOrderCountByCategory(dataset *reporting.OrderDataset) {
	fmt.Println("Order count by category")
	fmt.Println()

	data := make([]stat, 0)

	for _, c := range dataset.AllCategories() {
		data = append(data, stat{string(c), float64(dataset.NumOrdersByCategory(c))})
	}

	renderOrderCountByCategoryTable(data)
	renderOrderCountByCategoryGraph(data)
}

func renderOrderCountByCategoryTable(data []stat) {
	textData := make([][]string, 0)
	for _, d := range data {
		textData = append(textData, []string{d.x, fmt.Sprintf("%d", int(d.y))})
	}

	tap.Table(
		[]string{"Category", "Order count"},
		textData,
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderOrderCountByCategoryGraph(data []stat) {
	values := make([]barchart.BarData, 0)
	for _, stat := range data {
		values = append(
			values,
			barchart.BarData{
				Label:  stat.x,
				Values: []barchart.BarValue{{Name: "Order count", Value: stat.y, Style: blockStyle}}})
	}

	bc := barchart.New(
		160, 15,
		barchart.WithDataSet(values))

	// bc := barchart.New(
	// 	140, len(data)*2,
	// 	barchart.WithDataSet(values),
	// 	barchart.WithHorizontalBars())

	bc.Draw()

	fmt.Println(bc.View())
}
