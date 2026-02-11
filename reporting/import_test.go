package reporting_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"refurbed.com/hackathon/reporting"
)

func TestImportOrderDatasetFromCSV(t *testing.T) {
	in, err := os.Open("orders_v3.csv")
	require.NoError(t, err)
	t.Cleanup(func() { _ = in.Close() })

	dataset, err := reporting.ImportOrderDatasetFromCSV(in)
	require.NoError(t, err)

	t.Log("NumOrderItems", dataset.NumOrderItems())
	t.Log("NumOrders", dataset.NumOrders())
	t.Log("AOV", dataset.AOV())
	t.Log("Total revenue", dataset.TotalRevenue())

	for _, cat := range dataset.AllCategories() {
		t.Log("Orders By Category:", cat, dataset.NumOrdersByCategory(cat))
		t.Log("Return Rate By Category:", cat, dataset.ReturnRateByCategory(cat))
	}

}
