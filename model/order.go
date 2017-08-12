package model

import "time"

type Order struct {
	OrderID       string
	TransactionID string
	Sell          bool
	Buy           bool
	Date          time.Time
	Value         int
	Quantity      float64
	BuySuccess    func(*Event, *Order) `json:"-"`
	SellSuccess   func(*Event, *Order) `json:"-"`
}
