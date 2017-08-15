package bots

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/model"
)

type BotV1 struct {
	BotBase
	MinimumGapToBuy          float64
	MinimumQuantityToBuy     float64
	MaximumQuantityToBuy     float64
	MinimumPercentGapToSell  float64
	ScaleGapToBuy            string
	StrategyBuyPeriod        string
	StrategyBuyPercent       float64
	StrategySellPeriod       string
	StrategySellMetric       string
	StrategySellPercent      float64
	IsLookingUpwardVariation bool
	IsLowPassFilter          bool
	IsHighPassFilter         bool
}

func (b *BotV1) Update(history *history.History, event *model.Event, date time.Time) bool {
	isChanged := false

	b.Purge(date)

	b.MakeBuyOrder(history, event, date)

	b.MakeSellOrder(history, event, date)

	return isChanged
}

func (b *BotV1) Purge(date time.Time) {
	//log.Println(b.Orders)
	for orderID, order := range b.Orders {
		if date.Sub(order.Date) > time.Minute*15 {
			b.client.CancelOrder(order.OrderID)
			if b.Orders[orderID] == nil {
				log.Fatal("CAN'T FIND ORDER")
			}
			delete(b.Orders, orderID)
			isTxFound := 0
			for _, transaction := range b.Transactions {
				if transaction.TransactionID == order.TransactionID {
					transaction.InProgress = false
					isTxFound++
				}
			}
			if isTxFound != 1 {
				log.Println("TX NOT FOUND (PURGE) : ", isTxFound)
			}
		}
	}
}

func (b *BotV1) MakeBuyOrder(history *history.History, event *model.Event, date time.Time) bool {
	isChanged := false

	if !b.IsLookingUpwardVariation || b.IsPriceVariationGood(true, history, date) {
		priceToBuy := b.GetPriceToBuy(history, date)
		quantityToBuy := b.GetQuantityToBuy(priceToBuy, history, date)
		if b.IsPriceGapEnoughToBuy(priceToBuy, history, date) && quantityToBuy > 0.0 {
			if !b.IsLowPassFilter || priceToBuy <= b.LowPassFilter(history, date) {
				transaction := b.client.MakeBuyOrder(quantityToBuy, priceToBuy, priceToBuy, b.BuySuccess)
				if transaction != nil {
					transaction.Date = date
					b.Orders[transaction.OrderID] = transaction
					isChanged = true
					//log.Println("MAKE A BID : ", "VAL : ", priceToBuy, "QTY : ", quantityToBuy, "MIN : ", b.MinimumQuantityToBuy, "ORDERS : ", len(b.Orders))
				} else {
					log.Fatal("Cannot make order (BUY)")
				}
			} else {
				//log.Println(b.LowPassFilter(history, date), "/", priceToBuy)
			}
		}
	}

	return isChanged
}

func (b *BotV1) MakeSellOrder(history *history.History, event *model.Event, date time.Time) bool {
	isChanged := false

	if len(b.Transactions) > 0 {
		if (!b.IsLookingUpwardVariation) || b.IsPriceVariationGood(false, history, date) {
			priceToSell := b.GetPriceToSell(history, date)
			if priceToSell > 0 {
				for _, transaction := range b.Transactions {
					if priceToSell >= int(float64(transaction.TransactionValue)*b.MinimumPercentGapToSell) && !transaction.InProgress {
						if !b.IsHighPassFilter || priceToSell >= b.HighPassFilter(history, date) {
							order := b.client.MakeSellOrder(transaction.Quantity, priceToSell, transaction.OriginalValue, b.SellSuccess)
							if order != nil {
								order.TransactionID = transaction.TransactionID
								b.Orders[order.OrderID] = order
								transaction.InProgress = true
								isChanged = true
								//log.Println("MAKE A BID : ", "VAL : ", priceToSell, "QTY : ", transaction.Quantity)
							} else {
								log.Fatal("Cannot make order (SELL)")
							}
						}
					} else {
						//log.Println(priceToSell - int(float64(transaction.TransactionValue)*b.MinimumPercentGapToSell))
					}
				}
			}
		}
	}

	return isChanged
}

func (b *BotV1) IsOrderOnTransaction(transactionID string) bool {
	for _, order := range b.Orders {
		if order.TransactionID == transactionID {
			return true
		}
	}
	return false
}

func (b *BotV1) IsPriceGapEnoughToBuy(price int, history *history.History, date time.Time) bool {
	if len(b.Orders) == 0 && len(b.Transactions) == 0 {
		return true
	}

	for _, transaction := range b.Transactions {
		if transaction.Date.Sub(date) < time.Minute*15 {
			return false
		}

		if math.Abs(float64(transaction.TransactionValue)-float64(price))/float64(transaction.TransactionValue) <= b.MinimumGapToBuy || math.Abs(float64(transaction.OriginalValue)-float64(price))/float64(transaction.OriginalValue) <= b.MinimumGapToBuy {
			var stats *model.Statistic

			//log.Println("FIRST FILTER : ", math.Abs(float64(transaction.TransactionValue)-float64(price))/float64(transaction.TransactionValue), "/", b.MinimumGapToBuy)

			// Switch between strategies
			switch b.StrategyBuyPeriod {
			case "hour":
				stats = history.GetLastHourEvents(date)
				break
			case "six-hours":
				stats = history.GetLastSixHourEvents(date)
				break
			case "day":
				stats = history.GetLastDayEvents(date)
				break
			case "week":
				stats = history.GetLastWeekEvents(date)
				break
			case "month":
				stats = history.GetLastMonthEvents(date)
				break
			}

			if math.Abs(float64(stats.Max)-float64(price))/float64(price) < b.MinimumGapToBuy*1.25 {
				//log.Println("SECOND FILTER : ", math.Abs(float64(stats.Max)-float64(price))/float64(price), "/", b.MinimumGapToBuy*2)
				return false
			}
		}
	}

	for _, transaction := range b.Orders {
		if math.Abs(float64(transaction.OriginalValue-price)) <= b.MinimumGapToBuy {
			if transaction.Date.Sub(date) < time.Minute*15 {
				return false
			}

			var stats *model.Statistic

			//log.Println("FIRST FILTER : ", math.Abs(float64(transaction.OriginalValue-price)), "/", b.MinimumGapToBuy)

			// Switch between strategies
			switch b.StrategyBuyPeriod {
			case "hour":
				stats = history.GetLastHourEvents(date)
				break
			case "six-hours":
				stats = history.GetLastSixHourEvents(date)
				break
			case "day":
				stats = history.GetLastDayEvents(date)
				break
			case "week":
				stats = history.GetLastWeekEvents(date)
				break
			case "month":
				stats = history.GetLastMonthEvents(date)
				break
			}

			if math.Abs(float64(stats.Max)-float64(price))/float64(price) < b.MinimumGapToBuy*1.25 {
				//log.Println("SECOND FILTER : ", math.Abs(float64(stats.Max)-float64(price))/float64(price), "/", b.MinimumGapToBuy*2)
				return false
			}
		}
	}

	return true
}

func (b *BotV1) IsPriceVariationGood(buy bool, history *history.History, date time.Time) bool {
	var current *model.Statistic
	var previous *model.Statistic

	// Switch between strategies
	switch b.StrategyBuyPeriod {
	case "hour":
		current = history.GetLastHourEvents(date)
		previous = history.GetPreviousHourEvents(date)
		break
	case "six-hours":
		current = history.GetLastSixHourEvents(date)
		previous = history.GetPreviousSixHourEvents(date)
		break
	case "day":
		current = history.GetLastDayEvents(date)
		previous = history.GetPreviousDayEvents(date)
		break
	case "week":
		current = history.GetLastWeekEvents(date)
		previous = history.GetPreviousWeekEvents(date)
		break
	case "month":
		current = history.GetLastMonthEvents(date)
		previous = history.GetPreviousMonthEvents(date)
		break
	}

	if buy {
		if !previous.UpwardVariation && current.UpwardVariation {
			return true
		}
		return false
	} else {
		if previous.UpwardVariation && !current.UpwardVariation {
			return true
		}
		return false
	}
}

func (b *BotV1) GetQuantityToBuy(price int, history *history.History, date time.Time) float64 {
	quantity := b.MinimumQuantityToBuy / 1.5
	if price < history.GetLastSixHourEvents(date).Value {
		quantity = quantity * 1.25
	}
	if price < history.GetLastDayEvents(date).Value {
		quantity = quantity * 1.5
	}
	if price < history.GetLastWeekEvents(date).Value {
		quantity = quantity * 2.0
	}
	if price < history.GetLastMonthEvents(date).Value {
		quantity = quantity * 2.5
	}

	if float64(price)*quantity > float64(b.GetWalletValue()) {
		quantity = float64(float64(b.GetWalletValue()) / float64(price))
	}

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

func (b *BotV1) GetPriceToBuy(history *history.History, date time.Time) int {
	var stats *model.Statistic

	// Switch between strategies
	switch b.StrategyBuyPeriod {
	case "hour":
		stats = history.GetLastHourEvents(date)
		break
	case "six-hours":
		stats = history.GetLastSixHourEvents(date)
		break
	case "day":
		stats = history.GetLastDayEvents(date)
		break
	case "week":
		stats = history.GetLastWeekEvents(date)
		break
	case "month":
		stats = history.GetLastMonthEvents(date)
		break
	}

	return int(float64(stats.Min) * b.StrategyBuyPercent)
}

func (b *BotV1) GetPriceToSell(history *history.History, date time.Time) int {
	var stats *model.Statistic

	// Switch between strategies
	switch b.StrategySellPeriod {
	case "hour":
		stats = history.GetLastHourEvents(date)
		break
	case "six-hours":
		stats = history.GetLastSixHourEvents(date)
		break
	case "day":
		stats = history.GetLastDayEvents(date)
		break
	case "week":
		stats = history.GetLastWeekEvents(date)
		break
	case "month":
		stats = history.GetLastMonthEvents(date)
		break
	}

	switch b.StrategySellMetric {
	case "min":
		return int(float64(stats.Min) * b.StrategySellPercent)
	case "max":
		return int(float64(stats.Max) * b.StrategySellPercent)
	case "average":
		return int(float64(stats.Value) * b.StrategySellPercent)
	default:
		return 0
	}
}

func (b *BotV1) Display() string {
	return b.ID + " : " + b.StrategyBuyPeriod + ", " + b.StrategySellPeriod + ", Looking for variation : " + strconv.FormatBool(b.IsLookingUpwardVariation) + ", MinimumPercentGapToSell : " + strconv.FormatFloat(b.MinimumPercentGapToSell, 'f', -1, 64) + ", StrategyBuyPercent : " + strconv.FormatFloat(b.StrategyBuyPercent, 'f', -1, 64) + ", StrategySellPercent : " + strconv.FormatFloat(b.StrategySellPercent, 'f', -1, 64) + ", StrategySellMetric : " + b.StrategySellMetric + ", IsLowPassFilter : " + strconv.FormatBool(b.IsLowPassFilter) + ", IsHighPassFilter : " + strconv.FormatBool(b.IsHighPassFilter)
}

func (b *BotV1) LowPassFilter(history *history.History, date time.Time) int {
	var stats *model.Statistic

	// Switch between strategies
	switch b.StrategyBuyPeriod {
	case "hour":
		stats = history.GetLastHourEvents(date)
		break
	case "six-hours":
		stats = history.GetLastSixHourEvents(date)
		break
	case "day":
		stats = history.GetLastDayEvents(date)
		break
	case "week":
		stats = history.GetLastWeekEvents(date)
		break
	case "month":
		stats = history.GetLastMonthEvents(date)
		break
	}

	return int(float64(stats.Value) * 1.25)
}

func (b *BotV1) HighPassFilter(history *history.History, date time.Time) int {
	var stats *model.Statistic

	// Switch between strategies
	switch b.StrategySellPeriod {
	case "hour":
		stats = history.GetLastHourEvents(date)
		break
	case "six-hours":
		stats = history.GetLastSixHourEvents(date)
		break
	case "day":
		stats = history.GetLastDayEvents(date)
		break
	case "week":
		stats = history.GetLastWeekEvents(date)
		break
	case "month":
		stats = history.GetLastMonthEvents(date)
		break
	}

	return int(float64(stats.Value) * 0.75)
}
