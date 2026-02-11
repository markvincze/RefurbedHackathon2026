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

	t.Log(dataset.AOV())
}
