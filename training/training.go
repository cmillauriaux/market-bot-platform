package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/cmillauriaux/market-bot-platform/bots"
	"github.com/cmillauriaux/market-bot-platform/engine"
	"github.com/cmillauriaux/market-bot-platform/history"
	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/cmillauriaux/market-bot-platform/utils"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	end, err := time.Parse("2006-01-02", "2017-03-01")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Begining MBP training program")
	market := &market.TrainingClient{}
	log.Println("Market initialized")
	engine := engine.Init(market)
	engine.Bots, err = loadBots("./bots.json", market, engine.History)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Engine initialized")
	log.Println("Loading history...")
	engine.LoadHistory("../.coinbaseEUR.csv", end)
	log.Println("History loaded")
	go engine.LaunchSupervision("../supervision")
	log.Println("Begining training...")
	launchSimulation("../.coinbaseEUR.csv", end, engine.Bots, engine.History, market)
	log.Println("End of training")
	for {
		time.Sleep(time.Minute)
	}
}

func loadBots(filename string, market market.Market, statistics *history.History) ([]bots.Bot, error) {
	var b []*bots.BotV1
	currentBots := make([]bots.Bot, 0)

	file, _ := ioutil.ReadFile(filename)
	err := json.Unmarshal(file, &b)
	if err != nil {
		return nil, err
	}

	for _, bot := range b {
		bot.Init(market, statistics, 10000)
		currentBots = append(currentBots, bot)
	}

	return currentBots, nil
}

func savedBots(filename string, bots []bots.Bot) error {
	b, err := json.Marshal(bots)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func launchSimulation(filename string, start time.Time, botList []bots.Bot, history *history.History, market market.Market) {
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
	lastDay := time.Time{}
	var wg sync.WaitGroup
	for !transaction_closed {
		select {
		case event, transaction_ok := <-channel:
			wg.Wait()
			// Detect if CSV read is finish
			if !transaction_ok {
				transaction_closed = true
			} else {
				// Insert a new event in history
				if event.Date.After(start) {
					if event.Date.After(lastDay.Add(time.Hour * 24)) {
						log.Println(event.Date)
						stats := history.GetLastDayEvents(event.Date)
						log.Println("[Day] Open : ", stats.Open, "Close : ", stats.Close, " Average : ", stats.Value, " Min : ", stats.Min, " Max : ", stats.Max)
						week := history.GetLastWeekEvents(event.Date)
						log.Println("[Week] Open : ", week.Open, "Close : ", week.Close, " Average : ", week.Value, " Min : ", week.Min, " Max : ", week.Max)
						lastDay = event.Date
					}
					market.SimulateMarketTransaction(&event)
					history.InsertEvent(&event)
					for _, bot := range botList {
						wg.Add(1)
						go func(currentBot bots.Bot) {
							defer wg.Done()
							currentBot.Update(history, &event, event.Date)
						}(bot)
					}
				}
			}
		}
	}

	// Display performances informations
	log.Println("Training complete in ", counter.StopCount().Seconds(), "s")
}
