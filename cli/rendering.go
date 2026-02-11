package main

import "github.com/yarlson/tap"

func renderRevenueByDay() {
	renderRevenueByDayTable()
	renderRevenueByDayGraph()
}

func renderRevenueByDayTable() {
	tap.Table(
		[]string{"Day", "Revenue"},
		[][]string{
			[]string{"2026.01.01", "€ 120"},
			[]string{"2026.01.02", "€ 112"},
			[]string{"2026.01.03", "€ 270"},
			[]string{"2026.01.04", "€ 178"},
			[]string{"2026.01.05", "€ 153"},
			[]string{"2026.01.06", "€ 240"},
			[]string{"2026.01.07", "€ 225"},
		},
		tap.TableOptions{ShowBorders: true, HeaderStyle: tap.TableStyleBold})
}

func renderRevenueByDayGraph() {
}
