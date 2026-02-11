package reporting

import (
	"sync"
	"iter"
	"slices"
	"time"

	"github.com/shopspring/decimal"
)

type Category string

type OrderID string

type OrderItem struct {
	OrderID       OrderID
	OrderedAt     time.Time
	CustomerEmail string
	ItemName      string
	ItemSpecs     []ItemSpec
	ItemPrice     decimal.Decimal
	Commission    decimal.Decimal
	Refunded      decimal.Decimal
	PaymentStatus string
	Country       string
	ShippedAt     time.Time
	DeliveredAt   time.Time
	Category      []Category
}

type ItemSpec struct {
	Key      string
	RawValue string
}

type OrderDataset struct {
	allItems []OrderItem
	orders map[OrderID][]OrderItem

	aovOnce sync.Once
	aov     decimal.Decimal
}

func (ds *OrderDataset) AllItems() iter.Seq[OrderItem] {
	return slices.Values(ds.allItems)
}

func (ds *OrderDataset) AOV() decimal.Decimal {
	ds.aovOnce.Do(func() {
		total := decimal.Zero
		for item := range ds.AllItems() {
			total = total.Add(item.ItemPrice)
		}
		ds.aov = total.Div(decimal.NewFromInt(int64(len(ds.orders))))
	})
	return ds.aov
}
