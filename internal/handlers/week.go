package handlers

import (
	"cocon/internal/handlers/view"
	"cocon/internal/helpers"
	"net/http"
	"time"
)

const weekQueryLayout = "2006-01-02"

func requestedWeek(r *http.Request) time.Time {
	now := time.Now()
	weekStart := helpers.StartOfWeek(now)

	if raw := r.URL.Query().Get("week"); raw != "" {
		if date, err := time.ParseInLocation(weekQueryLayout, raw, now.Location()); err == nil {
			weekStart = helpers.StartOfWeek(date)
		}
	}
	return weekStart
}

func buildWeekNavigation(path string, weekStart time.Time) view.WeekNavigation {
	currentWeek := helpers.StartOfWeek(time.Now())
	query := "?week=" + weekStart.Format(weekQueryLayout)
	return view.WeekNavigation{
		Label:       helpers.WeekRange(weekStart),
		PreviousURL: path + "?week=" + weekStart.AddDate(0, 0, -7).Format(weekQueryLayout),
		CurrentURL:  path,
		NextURL:     path + "?week=" + weekStart.AddDate(0, 0, 7).Format(weekQueryLayout),
		Query:       query,
		IsCurrent:   weekStart.Equal(currentWeek),
	}
}
