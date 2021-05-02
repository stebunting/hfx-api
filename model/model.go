package model

import (
	"time"
)

type Currency struct {
	Code   string `json:"code" pg:",pk"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Exchange struct {
	Date     time.Time `json:"date" pg:",pk"`
	FromCode string    `json:"fromCode" pg:",pk"`
	ToCode   string    `json:"toCode" pg:",pk"`
	Rate     float64   `json:"rate"`
}
