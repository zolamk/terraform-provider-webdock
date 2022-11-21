package api

// Price model
type Price struct {
	// Price amount
	Amount int64 `json:"amount,omitempty"`

	// Price currency
	Currency string `json:"currency,omitempty"`
}

func errorStatus(code int) bool {
	return code < 200 || code > 299
}
