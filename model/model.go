package model

import (
	"time"
)

type Currency struct {
	Code   string `pg:",pk"`
	Name   string
	Symbol string
}

type Exchange struct {
	Date     time.Time `pg:",pk"`
	FromCode string    `pg:",pk"`
	ToCode   string    `pg:",pk"`
	Rate     float64
}
