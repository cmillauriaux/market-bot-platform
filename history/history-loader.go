package history

import (
	"log"

	"git.icysoft.fr/cedric/kraken-bot/model"
	"git.icysoft.fr/cedric/kraken-bot/utils"
)

func LoadHistory(filename string) (*History, error) {
	counter := utils.Counter{}
	counter.StartCount()
	channel := make(chan model.Event)
	history, err := InitHistory()

	if err != nil {
		return nil, err
	}

	go func() {
		utils.ReadCsv(filename, channel, 0)
	}()

	for {
		transaction_closed := false

		select {
		case event, transaction_ok := <-channel:
			if !transaction_ok {
				transaction_closed = true
			} else {
				err = history.InsertEvent(&event)
				if err != nil {
					return nil, err
				}
			}
		}

		if transaction_closed {
			break
		}
	}
	log.Println("Load History complete in ", counter.StopCount().Seconds(), "s")

	history.RefreshStatistics()

	return history, nil
}
