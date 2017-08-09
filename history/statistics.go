package history

import (
	"log"
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/cmillauriaux/market-bot-platform/utils"
	"github.com/jinzhu/now"
)

// Méthode pour la vue
func (h *History) InstantStatistics() *model.Statistic {
	// Register performance metrics
	counter := utils.Counter{}
	counter.StartCount()

	// Init datas
	events := make([]*model.Event, 0)
	it := h.Realtime.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)
		events = append(events, event)
	}
	statistic := h.ComputeStatistics(events)
	// Display performance metrics
	log.Println("Generate instant statistics in", counter.StopCount().Seconds(), "s")
	return statistic
}

func (h *History) LastHourEvents() *model.Statistic {
	// Register performance metrics
	counter := utils.Counter{}
	counter.StartCount()

	events := make([]*model.Event, 0)
	it := h.Realtime.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)
		if event.Date.After(time.Now().Add(-time.Hour)) {
			events = append(events, event)
		}
	}
	// Display performance metrics
	log.Println("Generate last hour statistics in", counter.StopCount().Seconds(), "s")

	statistic := h.ComputeStatistics(events)
	return statistic
}

func (h *History) YearsStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(YEAR, true, time.Unix(0, 0), time.Unix(0, 0)))
}

func (h *History) MonthsStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(MONTH, true, time.Unix(0, 0), time.Unix(0, 0)))
}

func (h *History) WeeksStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(WEEK, true, time.Unix(0, 0), time.Unix(0, 0)))
}

func (h *History) Last30DaysStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(DAY, true, time.Now().Add(-time.Hour*24*30).Truncate(time.Hour*24), time.Now()))
}

func (h *History) getStatistics(r Range, slicing bool, start time.Time, end time.Time) []*model.Statistic {
	statistics := make([]*model.Statistic, 0)
	currentHistory := make([]*model.Statistic, 0)
	beginDate := time.Unix(0, 0)
	it := h.Days.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Statistic)
		event.DisplayDate = event.Date.Format("2006-01-02")

		if (!start.Equal(time.Unix(0, 0)) && event.Date.Before(start)) || (!end.Equal(time.Unix(0, 0)) && event.Date.After(end)) {
			continue
		}
		// First event ever read : initialize the oldest event timestamp
		if beginDate == time.Unix(0, 0) {
			beginDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, compute statistics for the day and purge realtime
		if h.isANewPeriod(beginDate, event.Date, r, false) {
			// Compute and register statistics for the day
			statistic := h.AggregateStatistics(currentHistory)
			statistic.Date = beginDate
			statistic.DateFin = event.Date
			statistic.DisplayDate = event.Date.Format("2006-01-02")
			if statistic.Value > 0 {
				statistics = append(statistics, statistic)
			}

			// Init list
			beginDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Statistic, 0)

			// Remove from begin to current
			h.removeRealTimeUntil(h.currentDate)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
	}

	// Current period
	statistic := h.AggregateStatistics(currentHistory)
	statistic.Date = beginDate
	statistic.DateFin = time.Now()
	statistic.DisplayDate = statistic.DateFin.Format("2006-01-02")
	statistic.Partial = true
	statistics = append(statistics, statistic)

	// Return statistics
	return statistics
}

func (h *History) isANewPeriod(beginDate time.Time, currentDate time.Time, r Range, slicing bool) bool {
	if !slicing {
		now.FirstDayMonday = true
		beginDateNow := now.New(beginDate)
		switch r {
		case YEAR:
			if currentDate.After(beginDateNow.EndOfYear()) {
				return true
			}
			break
		case QUARTER:
			if currentDate.After(beginDateNow.EndOfQuarter()) {
				return true
			}
			break
		case MONTH:
			if currentDate.After(beginDateNow.EndOfMonth()) {
				return true
			}
			break
		case WEEK:
			if currentDate.After(beginDateNow.EndOfWeek()) {
				return true
			}
			break
		case DAY:
			if currentDate.After(beginDateNow.EndOfDay()) {
				return true
			}
			break
		default:
			return false
		}
	} else {
		switch r {
		case YEAR:
			if currentDate.After(beginDate.Add(time.Hour * 24 * 365)) {
				return true
			}
			break
		case QUARTER:
			if currentDate.After(beginDate.Add(time.Hour * 24 * 90)) {
				return true
			}
			break
		case MONTH:
			if currentDate.After(beginDate.Add(time.Hour * 24 * 30)) {
				return true
			}
			break
		case WEEK:
			if currentDate.After(beginDate.Add(time.Hour * 24 * 7)) {
				return true
			}
			break
		case DAY:
			if currentDate.After(beginDate.Add(time.Hour * 24)) {
				return true
			}
			break
		default:
			return false
		}
	}

	return false
}

func (h *History) MakeStatistics(stats []*model.Statistic) *model.Statistics {
	return &model.Statistics{Details: stats, Summary: h.AggregateStatistics(stats)}
}
