package entities

import "time"

type CurrencyRate struct {
	Currency string     `json:"currency"`
	Rate     float64    `json:"rate"`
	Date     *time.Time `json:"date"`
}
