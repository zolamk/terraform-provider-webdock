package api

// Defines values for PriceDTOCurrency.
const (
	DKK PriceDTOCurrency = "DKK"
	EUR PriceDTOCurrency = "EUR"
	USD PriceDTOCurrency = "USD"
)

// Price currency
type PriceDTOCurrency string

// Price model
type PriceDTO struct {
	// Price amount
	Amount *int64 `json:"amount,omitempty"`

	// Price currency
	Currency *PriceDTOCurrency `json:"currency,omitempty"`
}

func errorStatus(code int) bool {
	return code < 200 || code > 299
}
