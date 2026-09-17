package handlers

import (
	"net/http"
	"time"
)

func (app *Application) plannedMeals(w http.ResponseWriter, r *http.Request) {
	weekStart := startOfWeek(time.Now())

	mealsView, err := app.getPlannedMealsView(weekStart)
	if err != nil {
		app.serverError(w, err)
		return
	}

	groceries, err := app.getGroceriesView()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "planned_meals.tmpl", map[string]any{
		"WeekStart": mealsView.WeekStart,
		"Days":      mealsView.Days,
		"Groceries": groceries,
	})
}
