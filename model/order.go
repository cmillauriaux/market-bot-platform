package model

import "time"

type Order struct {
	OrderID          string
	TransactionID    string
	Sell             bool
	Buy              bool
	Date             time.Time
	DisplayDate      string
	OriginalValue    int
	OrderValue       int
	TransactionValue int
	Quantity         float64
	InProgress       bool
	BuySuccess       func(*Event, *Order) `json:"-"`
	SellSuccess      func(*Event, *Order) `json:"-"`
}

func (e *Order) GetOriginalValue() int {
	return int(e.OriginalValue*int(e.Quantity*100)) / 100
}

func (e *Order) GetOrderValue() int {
	return int(e.OrderValue*int(e.Quantity*100)) / 100
}

func (e *Order) GetTransactionValue() int {
	return int(e.TransactionValue*int(e.Quantity*100)) / 100
}

func (e *Order) GetPlusValue() int {
	return e.GetTransactionValue() - e.GetOriginalValue()
}
