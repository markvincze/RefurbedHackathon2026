package reporting

import (
	"iter"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/shopspring/decimal"
)

type Category string

type OrderID string

type OrderItem struct {
	OrderID        OrderID
	NumericOrderID int32
	OrderedAt      time.Time
	CustomerEmail  string
	ItemName       string
	ItemSpecs      []ItemSpec
	ItemPrice      decimal.Decimal
	Commission     decimal.Decimal
	Refunded       decimal.Decimal
	PaymentStatus  string
	Country        string
	ShippedAt      time.Time
	DeliveredAt    time.Time
	Category       []Category
}

type orderItemID = int32

type ItemSpec struct {
	Key      string
	RawValue string
}

type OrderDataset struct {
	allItems   []OrderItem
	orders     map[OrderID][]OrderItem
	features   features
	categories map[Category]struct{}

	aovOnce sync.Once
	aov     decimal.Decimal
}

type features struct {
	orderCategory     map[Category]*roaring.Bitmap
	orderItemCategory map[Category]*roaring.Bitmap
	returned          *roaring.Bitmap
}

type Order []OrderItem

func (ds *OrderDataset) AllItems() iter.Seq[OrderItem] {
	return slices.Values(ds.allItems)
}

func (ds *OrderDataset) AllOrders() iter.Seq[Order] {
	return func(yield func(Order) bool) {
		for _, items := range ds.orders {
			if !yield(items) {
				return
			}
		}
	}
}

func (ds *OrderDataset) add(item OrderItem) {
	itemID := orderItemID(len(ds.allItems))
	ds.allItems = append(ds.allItems, item)
	ds.orders[item.OrderID] = append(ds.orders[item.OrderID], item)

	if ds.features.returned == nil {
		ds.features.returned = roaring.New()
	}
	if ds.features.orderCategory == nil {
		ds.features.orderCategory = map[Category]*roaring.Bitmap{}
	}
	if ds.features.orderItemCategory == nil {
		ds.features.orderItemCategory = map[Category]*roaring.Bitmap{}
	}
	if ds.categories == nil {
		ds.categories = map[Category]struct{}{}
	}

	if !item.Refunded.IsZero() {
		ds.features.returned.Add(uint32(itemID))
	}
	if ds.features.orderCategory == nil {
		ds.features.orderCategory = map[Category]*roaring.Bitmap{}
	}
	if ds.features.orderItemCategory == nil {
		ds.features.orderItemCategory = map[Category]*roaring.Bitmap{}
	}
	for _, cat := range item.Category {
		ds.categories[cat] = struct{}{}
		orderCategoryBitmap := ds.features.orderCategory[cat]
		if orderCategoryBitmap == nil {
			orderCategoryBitmap = roaring.New()
			ds.features.orderCategory[cat] = orderCategoryBitmap
		}
		orderCategoryBitmap.Add(uint32(item.NumericOrderID))

		orderItemCategoryBitmap := ds.features.orderItemCategory[cat]
		if orderItemCategoryBitmap == nil {
			orderItemCategoryBitmap = roaring.New()
			ds.features.orderItemCategory[cat] = orderItemCategoryBitmap
		}
		orderItemCategoryBitmap.Add(uint32(itemID))
	}
}

func (ds *OrderDataset) AllCategories() []Category {
	all := slices.Collect(maps.Keys(ds.categories))
	slices.Sort(all)
	return all
}

func (ds *OrderDataset) NumOrdersByCategory(cat Category) int {
	bitmap := ds.features.orderCategory[cat]
	return int(bitmap.GetCardinality())
}

func (ds *OrderDataset) ReturnRateByCategory(cat Category) float64 {
	allInCategory := ds.features.orderItemCategory[cat]
	returnedInCategory := allInCategory.Clone()
	returnedInCategory.And(ds.features.returned)
	return float64(returnedInCategory.GetCardinality()) / float64(allInCategory.GetCardinality())
}

func (ds *OrderDataset) NumOrderItems() int {
	return len(ds.allItems)
}

func (ds *OrderDataset) NumOrders() int {
	return len(ds.orders)
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
