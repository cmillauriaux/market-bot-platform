package main

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/bots"
	"github.com/cmillauriaux/market-bot-platform/engine"
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/cmillauriaux/market-bot-platform/utils"
)

func main() {
	end, err := time.Parse("2006-01-02", "2016-08-01")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Begining MBP training program")
	market := &market.TrainingClient{}
	log.Println("Market initialized")
	engine := engine.Init(market)
	engine.Bots = loadBots(market)
	log.Println("Engine initialized")
	log.Println("Loading history...")
	engine.LoadHistory("../.coinbaseEUR.csv", end)
	log.Println("History loaded")
	log.Println("Begining training...")
	launchSimulation("../.coinbaseEUR.csv", end, engine.Bots, engine.History, market)
	log.Println("End of training")
	//engine.LaunchSupervision()
}

func loadBots(market market.Market) []bots.Bot {
	b := make([]bots.Bot, 0)
	bot := &bots.BotV1{MinimumGapToBuy: 100, MaximumQuantityToBuy: 0.01, MinimumQuantityToBuy: 0.01}
	bot.Init(market, 10000)
	b = append(b, bot)
	return b
}

func launchSimulation(filename string, start time.Time, bots []bots.Bot, history *history.History, market market.Market) {
	// Counter to measure performances
	counter := utils.Counter{}
	counter.StartCount()

	// Create a channel to make an async treatment
	channel := make(chan model.Event)

	// Launch CSV Reader
	go func() {
		err := utils.ReadCsv(filename, channel, 0)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Loop while there is lines to read in CSV
	transaction_closed := false
	for !transaction_closed {
		select {
		case event, transaction_ok := <-channel:
			// Detect if CSV read is finish
			if !transaction_ok {
				transaction_closed = true
			} else {
				// Insert a new event in history
				if event.Date.After(start) {
					market.SimulateMarketTransaction(&event)
					history.InsertEvent(&event)
					for _, bot := range bots {
						bot.Update(history, &event, event.Date)
					}
				}
			}
		}
	}

	// Display performances informations
	log.Println("Training complete in ", counter.StopCount().Seconds(), "s")
}
