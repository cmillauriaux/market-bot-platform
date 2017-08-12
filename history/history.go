package history

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/utils"
	"github.com/emirpasic/gods/maps/treemap"
	god_utils "github.com/emirpasic/gods/utils"

	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
)

type Range int

const (
	YEAR Range = iota
	QUARTER
	MONTH
	WEEK
	DAY
	SIX_HOURS
	THREE_HOURS
	HOUR
	FITEEN_MINUTES
	FIVE_MINUTES
	MINUTE
)

type History struct {
	Realtime    *treemap.Map
	currentDate time.Time
	Days        *treemap.Map
}

// InitHistory initialize a new history with all statistics maps ready to use
func InitHistory() *History {
	realtime := treemap.NewWith(god_utils.Int64Comparator)
	Days := treemap.NewWith(god_utils.Int64Comparator)
	return &History{
		Realtime: realtime,
		Days:     Days}
}

// InserEvent inserts a new event in history an refresh statistics if it's relevant
func (h *History) InsertEvent(event *model.Event) error {
	// For performance reasons, the key is the event timestamp
	key := event.Date.UnixNano()
	// Add event in realtime map : only statistics are saved in days, weeks, quarters, etc. maps
	h.Realtime.Put(key, event)
	// If the gap between oldest event and the event is at least a day, refresh statistics
	if event.Date.Sub(h.currentDate) > h.getRangeStep(DAY) {
		h.ComputeRealTime()
	}
	return nil
}

func (h *History) CompleteHistory(market market.Market, step time.Duration) error {
	counter := utils.Counter{}
	counter.StartCount()
	currentDate := h.currentDate

	for currentDate.Before(time.Now().Truncate(time.Hour * 24)) {
		statistic, err := market.GetStatistic(currentDate, currentDate.Add(time.Hour*24))
		statistic.Date = currentDate
		statistic.DateFin = currentDate.Add(time.Hour * 24)
		if err != nil {
			return err
		}

		if statistic.Value > 0 {
			h.Days.Put(statistic.Date.UnixNano(), statistic)
		}

		time.Sleep(market.GetSleepTimeBetweenRequests())
		currentDate = currentDate.Add(step)
	}

	log.Println("Complete History complete in ", counter.StopCount().Seconds(), "s")

	return nil
}

func (h *History) GetRealtimeInformations(market market.Market) error {
	counter := utils.Counter{}
	counter.StartCount()
	events, err := market.GetTransactions(time.Now().Add(-time.Hour*24).Truncate(time.Hour*24), time.Now())
	if err != nil {
		return err
	}
	for _, event := range events {
		h.Realtime.Put(event.Date.UnixNano(), event)
	}
	h.ComputeRealTime()
	log.Println("Load real time informations in ", counter.StopCount().Seconds(), "s")
	return nil
}

// ComputeRealTime look at every realtime event and make statistics if it's relevant
func (h *History) ComputeRealTime() {
	currentHistory := make([]*model.Event, 0)

	// For each event in realtime map
	it := h.Realtime.Iterator()
	for it.Next() {
		// Convert value in a event
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)

		// First event ever read : initialize the oldest event timestamp
		if h.currentDate == time.Unix(0, 0) {
			h.currentDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, compute statistics for the day and purge realtime
		if event.Date.Sub(h.currentDate) > h.getRangeStep(DAY) {
			// Compute and register statistics for the day
			statistic := h.ComputeStatistics(currentHistory)
			statistic.Date = h.currentDate
			statistic.DateFin = event.Date
			if statistic.Value > 0 {
				h.Days.Put(h.currentDate.UnixNano(), statistic)
			}

			// Init list
			h.currentDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Event, 0)

			// Remove from begin to current
			h.removeRealTimeUntil(h.currentDate)
		}

		// Add current day to list
		if !event.Date.Before(h.currentDate) {
			currentHistory = append(currentHistory, event)
		}
	}
}

// removeRealTimeUntil remove every events before a date
func (h *History) removeRealTimeUntil(date time.Time) {
	it := h.Realtime.Iterator()
	for it.Next() {
		key, value := it.Key(), it.Value()
		event := value.(*model.Event)

		// If event is before the param date, remove it
		if event.Date.Before(date.Add(-time.Hour * 24)) {
			h.Realtime.Remove(key)
		}
	}
}

func (h *History) getRangeStep(r Range) time.Duration {
	switch r {
	case YEAR:
		return time.Hour * 24 * 365
	case QUARTER:
		return time.Hour * 24 * 90
	case MONTH:
		return time.Hour * 24 * 30
	case WEEK:
		return time.Hour * 24 * 7
	case DAY:
		return time.Hour * 24 * 1
	default:
		return time.Hour * 24 * 1
	}
}

func (h *History) ComputeAverageValue(events []*model.Event) (int, float64) {
	totalQuantity := 0.0
	totalValue := 0
	nbValues := 0
	for _, event := range events {
		totalQuantity += event.Quantity
		totalValue += event.Value
		nbValues++
	}
	if nbValues == 0 {
		return 0, 0
	}
	return totalValue / nbValues, totalQuantity
}

func (h *History) ComputeStatistics(events []*model.Event) *model.Statistic {
	totalQuantity := 0.0
	totalValue := 0
	nbValues := 0
	min := 0
	max := 0
	open := 0
	close := 0
	for _, event := range events {
		totalQuantity += event.Quantity
		totalValue += event.Value
		nbValues++
		if min == 0 {
			open = event.Value
		}
		if min == 0 || event.Value < min {
			min = event.Value
		}
		if max == 0 || event.Value > max {
			max = event.Value
		}
		close = event.Value
	}

	delta := 0.0
	if open > 0 && close > 0 {
		if close > open {
			delta = (float64(close-open) / float64(close) * 100)
		} else {
			delta = (float64(open-close) / float64(open) * 100)
		}
	}

	value := 0
	if nbValues > 0 && totalValue > 0 {
		value = totalValue / nbValues
	}

	upward := false
	if close > open {
		upward = true
	}

	return &model.Statistic{Min: min, Max: max, Quantity: totalQuantity, Value: value, Delta: delta, Open: open, Close: close, UpwardVariation: upward}
}

func (h *History) AggregateStatistics(events []*model.Statistic) *model.Statistic {
	totalQuantity := 0.0
	totalValue := 0
	nbValues := 0
	min := 0
	max := 0
	open := 0
	close := 0
	for _, event := range events {
		totalQuantity += event.Quantity
		totalValue += event.Value
		nbValues++
		if min == 0 {
			open = event.Value
		}
		if min == 0 || event.Value < min {
			min = event.Value
		}
		if max == 0 || event.Value > max {
			max = event.Value
		}
		close = event.Value
	}

	delta := 0.0
	if open > 0 && close > 0 {
		if close > open {
			delta = (float64(close-open) / float64(close) * 100)
		} else {
			delta = (float64(open-close) / float64(open) * 100)
		}
	}

	value := 0
	if nbValues > 0 && totalValue > 0 {
		value = totalValue / nbValues
	}

	upward := false
	if close > open {
		upward = true
	}

	return &model.Statistic{Min: min, Max: max, Quantity: totalQuantity, Value: value, Delta: delta, Open: open, Close: close, UpwardVariation: upward}
}
