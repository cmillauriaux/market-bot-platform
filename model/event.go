package model

import "time"

type Event struct {
	OrderID       string
	TransactionID string
	Sell          bool
	Buy           bool
	Date          time.Time
	DisplayDate   string
	Value         int
	Quantity      float64
}

func (e *Event) GetEventValue() int {
	return int(e.Value*int(e.Quantity*100)) / 100
}
