package handlers

import "time"

// startOfWeek renvoie le lundi (00:00) de la semaine contenant t.
func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // dimanche -> 7, pour calculer l'offset depuis lundi
		weekday = 7
	}
	offset := weekday - 1
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -offset)
}
