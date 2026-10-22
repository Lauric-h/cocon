package handlers

import (
	"cocon/internal/handlers/view"
	"net/http"
)

func (app *Application) dashboardPage(w http.ResponseWriter, r *http.Request) {
	weekStart := requestedWeek(r)
	dashboardView, err := view.BuildDashboard(app.Meals, app.Groceries, app.Tasks, app.Calendar, weekStart)
	if err != nil {
		app.serverError(w, err)
		return
	}
	navigation := buildWeekNavigation("/", weekStart)
	dashboardView.WeekNavigation = navigation
	dashboardView.WeekQuery = navigation.Query
	app.render(w, http.StatusOK, "dashboard.tmpl", dashboardView)
}
