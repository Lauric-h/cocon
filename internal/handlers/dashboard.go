package handlers

import (
	"cocon/internal/handlers/view"
	"net/http"
)

func (app *Application) dashboardPage(w http.ResponseWriter, r *http.Request) {
	dashboardView, err := view.BuildDashboard(app.Meals, app.Groceries, app.Tasks, app.Calendar)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "dashboard.tmpl", dashboardView)
}
