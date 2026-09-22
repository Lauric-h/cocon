package helpers

import "time"

var weekdayFull = map[time.Weekday]string{
	time.Monday:    "Lundi",
	time.Tuesday:   "Mardi",
	time.Wednesday: "Mercredi",
	time.Thursday:  "Jeudi",
	time.Friday:    "Vendredi",
	time.Saturday:  "Samedi",
	time.Sunday:    "Dimanche",
}

var weekdayShort = map[time.Weekday]string{
	time.Monday:    "Lun",
	time.Tuesday:   "Mar",
	time.Wednesday: "Mer",
	time.Thursday:  "Jeu",
	time.Friday:    "Ven",
	time.Saturday:  "Sam",
	time.Sunday:    "Dim",
}

var months = map[time.Month]string{
	time.January:   "janvier",
	time.February:  "février",
	time.March:     "mars",
	time.April:     "avril",
	time.May:       "mai",
	time.June:      "juin",
	time.July:      "juillet",
	time.August:    "août",
	time.September: "septembre",
	time.October:   "octobre",
	time.November:  "novembre",
	time.December:  "décembre",
}

func DayFull(t time.Time) string  { return weekdayFull[t.Weekday()] }
func DayShort(t time.Time) string { return weekdayShort[t.Weekday()] }

// DateLong renvoie "mercredi 17 septembre 2026".
func DateLong(t time.Time) string {
	return weekdayFull[t.Weekday()] + " " + t.Format("2") + " " + months[t.Month()] + " " + t.Format("2006")
}

// StartOfWeek renvoie le lundi (00:00) de la semaine contenant t.
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	offset := weekday - 1
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -offset)
}
