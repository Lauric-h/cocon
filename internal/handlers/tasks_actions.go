package handlers

import (
	"cocon/internal/models"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) createTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulaire invalide", http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	if title == "" {
		http.Error(w, "titre requis", http.StatusBadRequest)
		return
	}
	category := models.TaskCategory(r.PostForm.Get("category"))

	var dueDate *time.Time
	if raw := r.PostForm.Get("due_date"); raw != "" {
		parsed, err := time.Parse("2006-01-02", raw)
		if err != nil {
			http.Error(w, "date invalide", http.StatusBadRequest)
			return
		}
		dueDate = &parsed
	}

	task := models.Task{
		Title:        title,
		TaskCategory: category,
		DueDate:      dueDate,
		Source:       "manual",
	}
	if _, err := app.Tasks.Create(task); err != nil {
		app.serverError(w, err)
		return
	}

	app.renderWeekTasks(w, startOfWeek(time.Now()))
}

func (app *Application) toggleTaskOccurrence(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if err := app.Tasks.ToggleOccurrenceDone(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderWeekTasks(w, startOfWeek(time.Now()))
}

func (app *Application) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if err := app.Tasks.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderWeekTasks(w, startOfWeek(time.Now()))
}

// renderWeekTasks recharge et affiche le fragment liste — réutilisé par toutes les actions tâches.
func (app *Application) renderWeekTasks(w http.ResponseWriter, weekStart time.Time) {
	view, err := app.getWeekTasksView(weekStart)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "task-list.tmpl", view)
}
