package reporting

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type rawOrderItemRow struct {
	OrderID       string
	OrderedAt     string
	CustomerEmail string
	ItemName      string
	ItemSpecs     string
	ItemPrice     string
	Commission    string
	Refunded      string
	PaymentStatus string
	Country       string
	ShippedAt     string
	DeliveredAt   string
	Category      string
}

type csvFieldIndex int

type csvField int

const (
	csvFieldOrderID csvField = iota
	csvFieldOrderedAt
	csvFieldCustomerEmail
	csvFieldItemName
	csvFieldItemSpecs
	csvFieldItemPrice
	csvFieldCommission
	csvFieldRefunded
	csvFieldPaymentStatus
	csvFieldCountry
	csvFieldShippedAt
	csvFieldDeliveredAt
	csvFieldCategory
)

func (f csvField) String() string {
	switch f {
	case csvFieldOrderID:
		return "order_id"
	case csvFieldOrderedAt:
		return "ordered_at"
	case csvFieldCustomerEmail:
		return "customer_email"
	case csvFieldItemName:
		return "item_name"
	case csvFieldItemSpecs:
		return "item_specs"
	case csvFieldItemPrice:
		return "item_price"
	case csvFieldCommission:
		return "commission"
	case csvFieldRefunded:
		return "refunded"
	case csvFieldPaymentStatus:
		return "payment_status"
	case csvFieldCountry:
		return "country"
	case csvFieldShippedAt:
		return "shipped_at"
	case csvFieldDeliveredAt:
		return "delivered_at"
	case csvFieldCategory:
		return "category"
	default:
		return "UNKNOWN FIELD"
	}
}

func requiredFields() []csvField {
	return []csvField{
		csvFieldOrderID,
		csvFieldOrderedAt,
		csvFieldCustomerEmail,
		csvFieldItemName,
		csvFieldItemSpecs,
		csvFieldItemPrice,
		csvFieldCommission,
		csvFieldRefunded,
		csvFieldPaymentStatus,
		csvFieldCountry,
		csvFieldShippedAt,
		csvFieldDeliveredAt,
		csvFieldCategory,
	}
}

func lookupFieldIndices(headerFields []string) ([]csvFieldIndex, error) {
	indices := make([]csvFieldIndex, len(requiredFields()))
	for _, reqfield := range requiredFields() {
		idx := slices.Index(headerFields, reqfield.String())
		if idx == -1 {
			return nil, fmt.Errorf("missing required field %q in header %q", reqfield, headerFields)
		}
		indices[reqfield] = csvFieldIndex(idx)
	}
	return indices, nil
}



func ImportOrderDatasetFromCSV(r io.Reader) (*OrderDataset, error) {
	csvr := csv.NewReader(r)

	headerFields, err := csvr.Read()
	if err != nil {
		return nil, fmt.Errorf("unexpected I/O error before CSV header: %w", err)
	}
	fieldIndices, err := lookupFieldIndices(headerFields)
	if err != nil {
		return nil, err
	}

	ds := &OrderDataset{
		allItems: make([]OrderItem, 0, 300_000),
		orders: map[OrderID][]OrderItem{},
	}
	for {
		fields, err := csvr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("read CSV row: %w", err)
		}
		raw := rawOrderItemRow{
			OrderID:       fields[fieldIndices[csvFieldOrderID]],
			OrderedAt:     fields[fieldIndices[csvFieldOrderedAt]],
			CustomerEmail: fields[fieldIndices[csvFieldCustomerEmail]],
			ItemName:      fields[fieldIndices[csvFieldItemName]],
			ItemSpecs:     fields[fieldIndices[csvFieldItemSpecs]],
			ItemPrice:     fields[fieldIndices[csvFieldItemPrice]],
			Commission:    fields[fieldIndices[csvFieldCommission]],
			Refunded:      fields[fieldIndices[csvFieldRefunded]],
			PaymentStatus: fields[fieldIndices[csvFieldPaymentStatus]],
			Country:       fields[fieldIndices[csvFieldCountry]],
			ShippedAt:     fields[fieldIndices[csvFieldShippedAt]],
			DeliveredAt:   fields[fieldIndices[csvFieldDeliveredAt]],
			Category:      fields[fieldIndices[csvFieldCategory]],
		}
		orderItem, err := parseOrderItem(raw)
		if err != nil {
			return nil, fmt.Errorf("parse order: %w", err)
		}
		ds.allItems = append(ds.allItems, orderItem)
		ds.orders[orderItem.OrderID] = append(ds.orders[orderItem.OrderID], orderItem)
	}
	return ds, nil
}

func parseOrderItem(raw rawOrderItemRow) (OrderItem, error) {
	errorf := func(format string, args ...any) (OrderItem, error) {
		return OrderItem{}, fmt.Errorf(format, args...)
	}
	parsedOrderedAt, err := time.Parse(time.RFC3339, raw.OrderedAt)
	if err != nil {
		return errorf("parse ordered_at: %w", err)
	}
	parsedItemPrice, err := decimal.NewFromString(raw.ItemPrice)
	if err != nil {
		return errorf("parse item_price: %w", err)
	}
	parsedCommission, err := decimal.NewFromString(raw.Commission)
	if err != nil {
		return errorf("parse commission: %w", err)
	}
	parsedRefunded, err := decimal.NewFromString(raw.Refunded)
	if err != nil {
		return errorf("parse refunded: %w", err)
	}
	var (
		parsedShippedAt   time.Time
		parsedDeliveredAt time.Time
	)
	if raw.ShippedAt != "" {
		parsedShippedAt, err = time.Parse(time.RFC3339, raw.ShippedAt)
		if err != nil {
			return errorf("parse shipped_at: %w", err)
		}
	}
	if raw.DeliveredAt != "" {
		parsedDeliveredAt, err = time.Parse(time.RFC3339, raw.DeliveredAt)
		if err != nil {
			return errorf("parse delivered_at: %w", err)
		}
	}
	return OrderItem{
		OrderID:       OrderID(raw.OrderID),
		OrderedAt:     parsedOrderedAt,
		CustomerEmail: raw.CustomerEmail,
		ItemName:      raw.ItemName,
		ItemSpecs:     parseItemSpecs(raw.ItemSpecs),
		ItemPrice:     parsedItemPrice,
		Commission:    parsedCommission,
		Refunded:      parsedRefunded,
		PaymentStatus: raw.PaymentStatus,
		Country:       raw.Country,
		ShippedAt:     parsedShippedAt,
		DeliveredAt:   parsedDeliveredAt,
		Category:      parseCategoryPath(raw.Category),
	}, nil
}

func parseCategoryPath(raw string) []Category {
	var path []Category
	for segment := range strings.FieldsFuncSeq(raw, isCategoryPathSeparator) {
		if segment == "" {
			continue
		}
		path = append(path, Category(segment))
	}
	return path
}

func isCategoryPathSeparator(r rune) bool {
	return r == '>'
}

func parseItemSpecs(raw string) []ItemSpec {
	var specs []ItemSpec
	for spec := range strings.FieldsFuncSeq(raw, isSpecSeparator) {
		key, rawValue, found := strings.Cut(spec, "=")
		if !found {
			continue
		}
		specs = append(specs, ItemSpec{
			Key:      key,
			RawValue: rawValue,
		})
	}
	return specs
}

func isSpecSeparator(r rune) bool {
	return r == '|'
}
