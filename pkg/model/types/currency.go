// Package types models the a variety of custom application types
package types

// Currency is a model that stores an amount as cents
type Currency int64

// Float64 converts the cents to the monetary unit
func (c Currency) Float64() float64 {
	return float64(c) / 100
}

// NewCurrency creates a Currecy value from a float64
func NewCurrency(v float64) Currency {
	return Currency(v * 100)
}
