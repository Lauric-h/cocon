package handlers

import "time"

var frenchWeekdayFull = map[time.Weekday]string{
	time.Monday:    "Lundi",
	time.Tuesday:   "Mardi",
	time.Wednesday: "Mercredi",
	time.Thursday:  "Jeudi",
	time.Friday:    "Vendredi",
	time.Saturday:  "Samedi",
	time.Sunday:    "Dimanche",
}

func frenchDayFull(t time.Time) string {
	return frenchWeekdayFull[t.Weekday()]
}
