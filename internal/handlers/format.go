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

var frenchWeekdayShort = map[time.Weekday]string{
	time.Monday:    "Lun",
	time.Tuesday:   "Mar",
	time.Wednesday: "Mer",
	time.Thursday:  "Jeu",
	time.Friday:    "Ven",
	time.Saturday:  "Sam",
	time.Sunday:    "Dim",
}

var frenchMonths = map[time.Month]string{
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

func frenchDayFull(t time.Time) string  { return frenchWeekdayFull[t.Weekday()] }
func frenchDayShort(t time.Time) string { return frenchWeekdayShort[t.Weekday()] }

// frenchDateLong renvoie "mercredi 17 septembre 2026".
func frenchDateLong(t time.Time) string {
	return frenchWeekdayFull[t.Weekday()] + " " + t.Format("2") + " " + frenchMonths[t.Month()] + " " + t.Format("2006")
}
