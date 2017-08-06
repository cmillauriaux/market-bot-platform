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
}
