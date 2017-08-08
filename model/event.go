package model

import "time"

type Event struct {
	OrderID     string
	Date        time.Time
	DisplayDate string
	Value       int
	Quantity    float64
}
