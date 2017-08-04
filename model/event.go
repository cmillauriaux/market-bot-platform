package model

import "time"

type Event struct {
	OrderID  string
	Date     time.Time
	Value    int
	Quantity float64
}
