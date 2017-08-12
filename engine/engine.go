package engine

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/bots"
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/supervision"
)

type Engine struct {
	market  market.Market
	History *history.History
	Bots    []bots.Bot
}

func Init(market market.Market) *Engine {
	engine := &Engine{market: market, History: history.InitHistory(), Bots: make([]bots.Bot, 0)}
	return engine
}

func (e *Engine) LoadHistory(fileName string, end time.Time) {
	history, err := history.LoadHistory(fileName, end)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Completing history...")
	err = history.CompleteHistory(e.market, time.Hour*24)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Getting realtime history...")
	err = history.GetRealtimeInformations(e.market)
	if err != nil {
		log.Fatal(err)
	}

	e.History = history
}

func (e *Engine) ConnectToMarket() {
	e.market.SubscribeToFlux(e.History.InsertEvent)
}

func (e *Engine) LaunchSupervision() {
	supervision.RunServer(e.History)
}
