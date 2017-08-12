package bots

import (
	"time"

	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
)

type Bot interface {
	Init(client market.Market, wallet int)
	Update(history *history.History, event *model.Event, date time.Time) bool
	Display()
}
