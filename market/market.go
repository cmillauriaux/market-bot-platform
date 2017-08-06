package market

import (
	"time"

	"git.icysoft.fr/cedric/kraken-bot/model"
)

type BroadcastEvent func(*model.Event) error

type Market interface {
	GetStatistic(start time.Time, end time.Time) (*model.Statistic, error)
	GetSleepTimeBetweenRequests() time.Duration
	SubscribeToFlux(broacastFn BroadcastEvent)
	GetTransactions(start time.Time, end time.Time) ([]*model.Event, error)
}
