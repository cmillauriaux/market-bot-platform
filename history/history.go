package history

import (
	"log"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"

	"git.icysoft.fr/cedric/kraken-bot/model"
)

type Range int

const (
	YEAR Range = iota
	QUARTER
	MONTH
	WEEK
	DAY
)

type History struct {
	realtime    *treemap.Map
	currentDate time.Time
	days        *treemap.Map
	Days        *treemap.Map
	Weeks       *treemap.Map
	Months      *treemap.Map
	Quarters    *treemap.Map
	Years       *treemap.Map
}

func InitHistory() (*History, error) {
	realtime := treemap.NewWith(utils.Int64Comparator)
	days := treemap.NewWith(utils.Int64Comparator)
	Days := treemap.NewWith(utils.Int64Comparator)
	Weeks := treemap.NewWith(utils.Int64Comparator)
	Months := treemap.NewWith(utils.Int64Comparator)
	Quarters := treemap.NewWith(utils.Int64Comparator)
	Years := treemap.NewWith(utils.Int64Comparator)
	return &History{
		realtime: realtime,
		days:     days,
		Days:     Days,
		Weeks:    Weeks,
		Months:   Months,
		Quarters: Quarters,
		Years:    Years}, nil
}

func (h *History) InsertEvent(event *model.Event) error {
	key := event.Date.UnixNano()
	h.realtime.Put(key, event)
	if event.Date.Sub(h.currentDate) > h.getRangeStep(DAY) {
		h.ComputeRealTime()
	}
	return nil
}

func (h *History) ComputeRealTime() {
	currentHistory := make([]*model.Event, 0)

	it := h.realtime.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)
		// First iteration
		if h.currentDate == time.Unix(0, 0) {
			h.currentDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, register day and purge realtime
		if event.Date.Sub(h.currentDate) > h.getRangeStep(DAY) {
			// Register day
			h.days.Put(h.currentDate.UnixNano(), currentHistory)

			// Init list
			h.currentDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Event, 0)

			//log.Println("CHECK : ", h.currentDate.String(), "[", h.realtime.Size(), "]")

			// Remove from begin to current
			h.removeRealTimeUntil(h.currentDate)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
	}
}

func (h *History) removeRealTimeUntil(date time.Time) {
	it := h.realtime.Iterator()
	for it.Next() {
		key, value := it.Key(), it.Value()
		event := value.(*model.Event)

		if event.Date.Before(date) {
			h.realtime.Remove(key)
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

func (h *History) RefreshStatistics() {
	h.refreshStatistic(DAY)
	//h.refreshStatistic(WEEK)
	//h.refreshStatistic(MONTH)
	//h.refreshStatistic(QUARTER)
	//h.refreshStatistic(YEAR)
}

func (h *History) refreshStatistic(r Range) {
	currentDate := time.Unix(0, 0)
	lastDate := time.Unix(0, 0)
	currentHistory := make([]*model.Event, 0)

	it := h.days.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		events := value.([]*model.Event)

		for _, event := range events {
			// First iteration
			if currentDate == time.Unix(0, 0) {
				currentDate = event.Date.Truncate(time.Hour * 24)
			}

			// If day is different, register day
			if event.Date.Sub(currentDate) > h.getRangeStep(r) {
				// Register statistics
				statistic := h.ComputeStatistics(currentHistory)
				statistic.Date = currentDate
				statistic.DateFin = lastDate
				log.Println(statistic.Display())

				// Init list
				currentDate = event.Date.Truncate(time.Hour * 24)
				currentHistory = make([]*model.Event, 0)
			}

			// Add current day to list
			currentHistory = append(currentHistory, event)
			lastDate = event.Date
		}
	}
	// Force partial statistics for last events
	statistic := h.ComputeStatistics(currentHistory)
	statistic.Date = currentDate
	statistic.DateFin = lastDate
	log.Println(statistic.Display())
}

func (h *History) ComputeStatistics(events []*model.Event) *model.Statistic {
	totalQuantity := 0.0
	totalValue := 0
	nbValues := 0
	min := 0
	max := 0
	for _, event := range events {
		totalQuantity += event.Quantity
		totalValue += event.Value
		nbValues++
		if min == 0 || event.Value < min {
			min = event.Value
		}
		if max == 0 || event.Value > max {
			max = event.Value
		}
	}

	delta := 0.0
	if min > 0 {
		delta = (float64(max-min) / float64(min) * 100)
	}

	value := 0
	if nbValues > 0 {
		value = totalValue / nbValues
	}

	return &model.Statistic{Min: min, Max: max, Quantity: totalQuantity, Value: value, Delta: delta}
}
