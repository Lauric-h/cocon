package handlers

import (
	"cocon/internal/models"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) upsertMeal(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulaire invalide", http.StatusBadRequest)
		return
	}

	dateStr := r.PostForm.Get("date")
	mealSlot := models.MealSlot(r.PostForm.Get("meal_slot"))
	dish := r.PostForm.Get("dish")
	notes := r.PostForm.Get("notes")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "date invalide", http.StatusBadRequest)
		return
	}

	switch mealSlot {
	case models.MealSlotBreakfast, models.MealSlotNoon, models.MealSlotDinner:
		// valide
	default:
		http.Error(w, "type de repas invalide", http.StatusBadRequest)
		return
	}

	if err := app.Meals.Upsert(date, mealSlot, dish, notes); err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "meal-cell.tmpl", map[string]any{
		"Date":     date,
		"MealSlot": mealSlot,
		"Dish":     dish,
		"Notes":    notes,
	})
}

func (app *Application) deleteMeal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if err := app.Meals.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
