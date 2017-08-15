package history

import (
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
	"github.com/jinzhu/now"
)

// Méthode pour la vue
func (h *History) InstantStatistics() *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, time.Now().Add(-time.Hour*24), time.Now()))
}

func (h *History) LastHourEvents() *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, time.Now().Add(-time.Hour), time.Now()))
}

func (h *History) GetLastHourEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, date.Add(-time.Hour), date))
}

func (h *History) GetLastSixHourEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, date.Add(-time.Hour*6), date))
}

func (h *History) GetPreviousSixHourEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, date.Add(-time.Hour*12), date.Add(-time.Hour*6)))
}

func (h *History) GetLastDayEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(HOUR, true, date.Add(-time.Hour*24), date))
}

func (h *History) GetPreviousDayEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(HOUR, true, date.Add(-time.Hour*48), date.Add(-time.Hour*24)))
}

func (h *History) GetLastWeekEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatistics(SIX_HOURS, true, date.Add(-time.Hour*24*7), date))
}

func (h *History) GetPreviousWeekEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatistics(SIX_HOURS, true, date.Add(-time.Hour*24*14), date.Add(-time.Hour*24*7)))
}

func (h *History) GetLastMonthEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatistics(DAY, true, date.Add(-time.Hour*24*30), date))
}

func (h *History) GetPreviousMonthEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatistics(DAY, true, date.Add(-time.Hour*24*60), date.Add(-time.Hour*24*30)))
}

func (h *History) GetPreviousHourEvents(date time.Time) *model.Statistic {
	return h.AggregateStatistics(h.getStatisticsFromEvents(MINUTE, true, date.Add(-time.Hour*2), date.Add(-time.Hour)))
}

func (h *History) LastHourStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatisticsFromEvents(MINUTE, true, time.Now().Add(-time.Hour), time.Now()))
}

func (h *History) LastSixHoursStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatisticsFromEvents(FITEEN_MINUTES, true, time.Now().Add(-time.Hour*6), time.Now()))
}

func (h *History) LastDayStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatisticsFromEvents(HOUR, true, time.Now().Add(-time.Hour*24), time.Now()))
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

func (h *History) Last7DaysStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(HOUR, true, time.Now().Add(-time.Hour*24*7).Truncate(time.Hour*24), time.Now()))
}

func (h *History) ForeverStatistics() *model.Statistics {
	return h.MakeStatistics(h.getStatistics(DAY, true, time.Unix(0, 0), time.Unix(0, 0)))
}

func (h *History) getStatisticsFromEvents(r Range, slicing bool, start time.Time, end time.Time) []*model.Statistic {
	statistics := make([]*model.Statistic, 0)
	currentHistory := make([]*model.Event, 0)
	beginDate := time.Unix(0, 0)
	it := h.Realtime.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Event)

		if (!start.Equal(time.Unix(0, 0)) && event.Date.Before(start)) || (!end.Equal(time.Unix(0, 0)) && event.Date.After(end)) {
			continue
		}
		// First event ever read : initialize the oldest event timestamp
		if beginDate == time.Unix(0, 0) {
			beginDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, compute statistics for the day and purge realtime
		if h.isANewPeriod(beginDate, event.Date, r, slicing) {
			// Compute and register statistics for the day
			statistic := h.ComputeStatistics(currentHistory)
			statistic.Date = beginDate
			statistic.DateFin = event.Date
			if r == MINUTE || r == FIVE_MINUTES || r == FITEEN_MINUTES || r == HOUR || r == THREE_HOURS || r == SIX_HOURS {
				statistic.DisplayDate = event.Date.Format("2006-01-02 15:04:05")
			} else {
				statistic.DisplayDate = event.Date.Format("2006-01-02")
			}

			if statistic.Value > 0 {
				statistics = append(statistics, statistic)
			}

			// Init list
			beginDate = event.Date
			currentHistory = make([]*model.Event, 0)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
	}

	// Current period
	statistic := h.ComputeStatistics(currentHistory)
	statistic.Date = beginDate
	statistic.DateFin = time.Now()
	if r == MINUTE || r == FIVE_MINUTES || r == FITEEN_MINUTES || r == HOUR || r == THREE_HOURS || r == SIX_HOURS {
		statistic.DisplayDate = statistic.Date.Format("2006-01-02 15:04:05")
	} else {
		statistic.DisplayDate = statistic.Date.Format("2006-01-02")
	}
	statistic.Partial = true
	statistics = append(statistics, statistic)

	// Return statistics
	return statistics
}

func (h *History) getStatistics(r Range, slicing bool, start time.Time, end time.Time) []*model.Statistic {
	statistics := make([]*model.Statistic, 0)
	currentHistory := make([]*model.Statistic, 0)
	beginDate := time.Unix(0, 0)
	it := h.Days.Iterator()
	for it.Next() {
		_, value := it.Key(), it.Value()
		event := value.(*model.Statistic)

		if (!start.Equal(time.Unix(0, 0)) && event.Date.Before(start)) || (!end.Equal(time.Unix(0, 0)) && event.Date.After(end)) {
			continue
		}
		// First event ever read : initialize the oldest event timestamp
		if beginDate == time.Unix(0, 0) {
			beginDate = event.Date.Truncate(time.Hour * 24)
		}

		// If day is different, compute statistics for the day and purge realtime
		if h.isANewPeriod(beginDate, event.Date, r, slicing) {
			// Compute and register statistics for the day
			statistic := h.AggregateStatistics(currentHistory)
			statistic.Date = beginDate
			statistic.DateFin = event.Date
			if r == MINUTE || r == FIVE_MINUTES || r == FITEEN_MINUTES || r == HOUR || r == THREE_HOURS || r == SIX_HOURS {
				statistic.DisplayDate = statistic.Date.Format("2006-01-02 15:04:05")
			} else {
				statistic.DisplayDate = statistic.Date.Format("2006-01-02")
			}
			if statistic.Value > 0 {
				statistics = append(statistics, statistic)
			}

			// Init list
			beginDate = event.Date.Truncate(time.Hour * 24)
			currentHistory = make([]*model.Statistic, 0)
		}

		// Add current day to list
		currentHistory = append(currentHistory, event)
	}

	// Current period
	statistic := h.AggregateStatistics(currentHistory)
	statistic.Date = beginDate
	statistic.DateFin = time.Now()
	if r == MINUTE || r == FIVE_MINUTES || r == FITEEN_MINUTES || r == HOUR || r == THREE_HOURS || r == SIX_HOURS {
		statistic.DisplayDate = statistic.Date.Format("2006-01-02 15:04:05")
	} else {
		statistic.DisplayDate = statistic.Date.Format("2006-01-02")
	}
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
		case SIX_HOURS:
			if currentDate.After(beginDate.Add(time.Hour * 6)) {
				return true
			}
			break
		case THREE_HOURS:
			if currentDate.After(beginDate.Add(time.Hour * 3)) {
				return true
			}
			break
		case HOUR:
			if currentDate.After(beginDate.Add(time.Hour)) {
				return true
			}
			break
		case FITEEN_MINUTES:
			if currentDate.After(beginDate.Add(time.Minute * 15)) {
				return true
			}
			break
		case FIVE_MINUTES:
			if currentDate.After(beginDate.Add(time.Minute * 5)) {
				return true
			}
			break
		case MINUTE:
			if currentDate.After(beginDate.Add(time.Minute)) {
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
