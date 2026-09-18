package handlers

import "net/http"

func (app *Application) dashboardPage(w http.ResponseWriter, r *http.Request) {
	view, err := app.getDashboardView()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "dashboard.tmpl", view)
}
