package bots

import (
	"log"
	"time"

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

func (b *BotV1) Update(event *model.Event, date time.Time) bool {
	isChanged := false

	b.Purge(date)

	transaction := b.client.MakeBuyOrder(0.1, event.Value, b.BuySuccess)
	if transaction != nil {
		b.Orders[transaction.OrderID] = transaction
		isChanged = true
		//log.Println("MAKE A BID : ", "VAL : ", priceToBuy, "QTY : ", quantityToBuy, "MIN : ", b.MinimumQuantityToBuy)
	}

	return isChanged
}

func (b *BotV1) Purge(date time.Time) {
	for _, order := range b.Orders {
		if date.Sub(order.Date) > time.Minute*15 {
			b.client.CancelOrder(order.OrderID)
			delete(b.Orders, order.OrderID)
			delete(b.ReverseOrders, order.TransactionID)
			log.Println("PURGE")
		}
	}
}
