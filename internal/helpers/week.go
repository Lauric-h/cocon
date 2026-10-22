package helpers

import (
	"fmt"
	"time"
)

func WeekRange(start time.Time) string {
	end := start.AddDate(0, 0, 6)
	if start.Year() == end.Year() && start.Month() == end.Month() {
		return fmt.Sprintf("%d au %d %s %d", start.Day(), end.Day(), months[end.Month()], end.Year())
	}
	if start.Year() == end.Year() {
		return fmt.Sprintf("%s au %s %d", DayMonth(start), DayMonth(end), end.Year())
	}
	return fmt.Sprintf("%s %d au %s %d", DayMonth(start), start.Year(), DayMonth(end), end.Year())
}
