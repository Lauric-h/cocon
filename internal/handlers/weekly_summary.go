package handlers

import (
	"cocon/internal/handlers/view"
	"cocon/internal/helpers"
	"cocon/internal/pdfexport"
	"net/http"
	"time"
)

func (app *Application) downloadWeeklySummary(w http.ResponseWriter, r *http.Request) {
	weekStart := helpers.StartOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekRangeLabel := helpers.DayMonth(weekStart) + " au " + helpers.DayMonth(weekEnd)

	days, err := view.BuildWeeklySummaryDays(app.Meals, app.Calendar, weekStart)
	if err != nil {
		app.serverError(w, err)
		return
	}

	tasks, err := view.BuildWeekTasks(app.Tasks, weekStart)
	if err != nil {
		app.serverError(w, err)
		return
	}

	pdfBytes, err := pdfexport.BuildWeeklySummary(weekRangeLabel, days, tasks)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="cocon-semaine.pdf"`)
	w.Write(pdfBytes)
}
