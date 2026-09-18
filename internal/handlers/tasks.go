package handlers

import "net/http"

func (app *Application) tasksPage(w http.ResponseWriter, r *http.Request) {
	view, err := app.getAllTasksView()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "tasks.tmpl", view)
}
