package main

import (
	"github.com/cmillauriaux/market-bot-platform/engine"
	"github.com/cmillauriaux/market-bot-platform/market"
)

func main() {
	market := market.InitMarket()
	engine := engine.Init(market)
	engine.LoadHistory(".krakenEUR-lite.csv")
	engine.ConnectToMarket()
	engine.LaunchSupervision()
}
