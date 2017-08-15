package market

import (
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
)

type BroadcastEvent func(*model.Event) error

type Market interface {
	GetStatistic(start time.Time, end time.Time) (*model.Statistic, error)
	GetSleepTimeBetweenRequests() time.Duration
	SubscribeToFlux(broacastFn BroadcastEvent)
	GetTransactions(start time.Time, end time.Time) ([]*model.Event, error)
	CancelOrder(orderId string)
	SimulateMarketTransaction(transaction *model.Event)
	MakeBuyOrder(size float64, value int, originalValue int, callback func(*model.Event, *model.Order)) *model.Order
	MakeSellOrder(size float64, value int, originalValue int, callback func(*model.Event, *model.Order)) *model.Order
}
