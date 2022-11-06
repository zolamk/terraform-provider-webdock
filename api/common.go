package api

// Defines values for PriceDTOCurrency.
const (
	DKK Currency = "DKK"
	EUR Currency = "EUR"
	USD Currency = "USD"
)

// Price currency
type Currency string

// Price model
type Price struct {
	// Price amount
	Amount *int64 `json:"amount,omitempty"`

	// Price currency
	Currency *Currency `json:"currency,omitempty"`
}

func errorStatus(code int) bool {
	return code < 200 || code > 299
}
