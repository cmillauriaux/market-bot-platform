package history

import (
	"log"

	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/cmillauriaux/market-bot-platform/utils"
)

// LoadHistory load an history from a CSV file
func LoadHistory(filename string) (*History, error) {
	// Counter to measure performances
	counter := utils.Counter{}
	counter.StartCount()

	// Create a channel to make an async treatment
	channel := make(chan model.Event)

	// Init a new history
	history, err := InitHistory()
	if err != nil {
		return nil, err
	}

	// Launch CSV Reader
	go func() {
		utils.ReadCsv(filename, channel, 0)
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
				err = history.InsertEvent(&event)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// Display performances informations
	log.Println("Load History complete in ", counter.StopCount().Seconds(), "s")

	// Refresh statistics once data are loaded
	history.ComputeRealTime()

	return history, nil
}
