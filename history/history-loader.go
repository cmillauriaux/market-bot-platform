package history

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/cmillauriaux/market-bot-platform/utils"
)

// LoadHistory load an history from a CSV file
func (h *History) LoadHistory(filename string, end time.Time) error {
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
				if !end.IsZero() && event.Date.After(end) {
					log.Println("Load History until ", event.Date)
					return nil
				}
				// Insert a new event in history
				err := h.InsertEvent(&event)
				if err != nil {
					return err
				}
			}
		}
	}

	// Display performances informations
	log.Println("Load History complete in ", counter.StopCount().Seconds(), "s")

	// Refresh statistics once data are loaded
	h.ComputeRealTime()

	return nil
}
