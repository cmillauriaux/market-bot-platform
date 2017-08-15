package main

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/engine"
	"github.com/cmillauriaux/market-bot-platform/market"
)

func main() {
	log.Println("Starting MBP")
	market := market.InitMarket()
	log.Println("Market initialized")
	engine := engine.Init(market)
	log.Println("Engine initialized")
	log.Println("Loading history...")
	engine.LoadHistory(".coinbaseEUR.csv", time.Time{})
	log.Println("History loaded")
	engine.ConnectToMarket()
	log.Println("Connected to market")
	log.Println("Launching supervision...")
	engine.LaunchSupervision("./supervision")
}
