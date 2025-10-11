package util

var (
	USD = "USD"
	EUR = "EUR"
	KES = "KES"
)

//returns true if currency is supported
func IsSupported(currency string) bool {
	switch currency {
	case USD, EUR, KES:
		return true
	}
	return false
}