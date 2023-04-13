package util

const (
	// USD is the currency code for US Dollar
	USD = "USD"
	// EUR is the currency code for Euro
	EUR = "EUR"
	// CAD is the currency code for Canadian Dollar
	CAD = "CAD"
)

// IsSupportedCurrency checks if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}
