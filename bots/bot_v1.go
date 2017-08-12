package bots

import (
	"math"
	"time"

	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/model"
)

type BotV1 struct {
	BotBase
	MinimumGapToBuy         float64
	MinimumQuantityToBuy    float64
	MaximumQuantityToBuy    float64
	MinimumPercentGapToSell float64
	ScaleGapToBuy           string
	StrategyBuy             string
	StrategySell            string
}

func (b *BotV1) Update(history *history.History, event *model.Event, date time.Time) bool {
	isChanged := false

	b.Purge(date)

	b.MakeBuyOrder(history, event, date)

	return isChanged
}

func (b *BotV1) Purge(date time.Time) {
	for _, order := range b.Orders {
		if date.Sub(order.Date) > time.Minute*15 {
			b.client.CancelOrder(order.OrderID)
			delete(b.Orders, order.OrderID)
			delete(b.ReverseOrders, order.TransactionID)
		}
	}
}

func (b *BotV1) MakeBuyOrder(history *history.History, event *model.Event, date time.Time) bool {
	isChanged := false

	// Get current hour statistics
	lastHour := history.GetLastHourEvents(date)

	// Get previsous hour statistics
	previousHour := history.GetPreviousHourEvents(date)

	if !previousHour.UpwardVariation && lastHour.UpwardVariation {
		priceToBuy := lastHour.Min
		quantityToBuy := b.GetQuantityToBuy(priceToBuy)
		if b.IsPriceGapEnoughToBuy(priceToBuy) && quantityToBuy > 0.0 {
			transaction := b.client.MakeBuyOrder(quantityToBuy, priceToBuy, b.BuySuccess)
			if transaction != nil {
				b.Orders[transaction.OrderID] = transaction
				isChanged = true
				//log.Println("MAKE A BID : ", "VAL : ", priceToBuy, "QTY : ", quantityToBuy, "MIN : ", b.MinimumQuantityToBuy, "ORDERS : ", len(b.Orders))
			}
		}
	}

	return isChanged
}

func (b *BotV1) IsPriceGapEnoughToBuy(price int) bool {
	if len(b.Orders) == 0 {
		return true
	}

	for _, transaction := range b.Transactions {
		if math.Abs(float64(transaction.Value-price)) > b.MinimumGapToBuy {
			return true
		}
	}
	for _, transaction := range b.Orders {
		if math.Abs(float64(transaction.Value-price)) > b.MinimumGapToBuy {
			return true
		}
	}
	return false
}

func (b *BotV1) GetQuantityToBuy(price int) float64 {
	quantity := float64(float64(b.GetWalletValue()) / float64(price))
	if quantity > b.MaximumQuantityToBuy {
		return b.MaximumQuantityToBuy
	}
	if quantity < b.MinimumQuantityToBuy {
		return 0
	}
	if quantity < 0 {
		return 0
	}
	return quantity
}
