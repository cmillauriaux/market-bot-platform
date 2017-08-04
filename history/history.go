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
	db       *treemap.Map
	days     *treemap.Map
	Days     *treemap.Map
	Weeks    *treemap.Map
	Months   *treemap.Map
	Quarters *treemap.Map
	Years    *treemap.Map
}

func InitHistory() (*History, error) {
	db := treemap.NewWith(utils.Int64Comparator)
	days := treemap.NewWith(utils.Int64Comparator)
	Days := treemap.NewWith(utils.Int64Comparator)
	Weeks := treemap.NewWith(utils.Int64Comparator)
	Months := treemap.NewWith(utils.Int64Comparator)
	Quarters := treemap.NewWith(utils.Int64Comparator)
	Years := treemap.NewWith(utils.Int64Comparator)
	return &History{
		db:       db,
		days:     days,
		Days:     Days,
		Weeks:    Weeks,
		Months:   Months,
		Quarters: Quarters,
		Years:    Years}, nil
}

func (h *History) InsertEvent(event *model.Event) error {
	key := event.Date.UnixNano()
	h.db.Put(key, event)
	return nil
}

func (h *History) GetNbEvents() int64 {
	return int64(h.db.Size())
}

func (h *History) ComputeDays() {
	currentDate := time.Unix(0, 0)
	currentHistory := make([]*model.Event, 0)

	it := h.db.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)
		// First iteration
		if currentDate == time.Unix(0, 0) {
			currentDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, register day
		if event.Date.Sub(currentDate) > h.getRangeStep(DAY) {
			// Register day
			h.days.Put(currentDate.UnixNano(), currentHistory)

			// Init list
			currentDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Event, 0)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
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
	//h.refreshStatistic(DAY)
	//h.refreshStatistic(WEEK)
	h.refreshStatistic(MONTH)
	//h.refreshStatistic(QUARTER)
	//h.refreshStatistic(YEAR)
}

func (h *History) refreshStatistic(r Range) {
	currentDate := time.Unix(0, 0)
	currentHistory := make([]*model.Event, 0)

	it := h.db.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)
		// First iteration
		if currentDate == time.Unix(0, 0) {
			currentDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, register day
		if event.Date.Sub(currentDate) > h.getRangeStep(r) {
			// Register statistics
			statistic := h.ComputeStatistics(currentHistory)
			statistic.Date = currentDate
			log.Println(statistic.Display())

			// Init list
			currentDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Event, 0)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
	}
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

	return &model.Statistic{Min: min, Max: max, Quantity: totalQuantity, Value: totalValue / nbValues, Delta: (float64(max-min) / float64(min) * 100)}
}
