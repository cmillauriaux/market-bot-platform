package engine

import (
	"log"

	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/supervision"
)

type Engine struct {
	market  market.Market
	History *history.History
}

func Init(market market.Market) *Engine {
	engine := &Engine{market: market}

	return engine
}

func (e *Engine) LoadHistory(fileName string) {
	history, err := history.LoadHistory(".krakenEUR-lite.csv")
	if err != nil {
		log.Fatal(err)
	}
	/*err = history.CompleteHistory(e.market, time.Hour*24)
	if err != nil {
		log.Fatal(err)
	}*/
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
