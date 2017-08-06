package main

import (
	"log"
	"time"

	"git.icysoft.fr/cedric/kraken-bot/history"
	"git.icysoft.fr/cedric/kraken-bot/market"
)

func main() {

	history, err := history.LoadHistory(".krakenEUR-lite.csv")
	if err != nil {
		log.Fatal(err)
	}
	market := market.InitMarket()
	err = history.CompleteHistory(market, time.Hour*24)
	if err != nil {
		log.Fatal(err)
	}
	err = history.GetRealtimeInformations(market)
	if err != nil {
		log.Fatal(err)
	}
	market.SubscribeToFlux(history.InsertEvent)
	for {
		time.Sleep(time.Hour)
	}
}
