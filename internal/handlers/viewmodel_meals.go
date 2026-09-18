package handlers

import (
	"cocon/internal/models"
	"time"
)

var mealSlots = []models.MealSlot{
	models.MealSlotBreakfast,
	models.MealSlotNoon,
	models.MealSlotDinner,
}

type MealCellView struct {
	Date     time.Time
	MealSlot models.MealSlot
	Dish     string
	Notes    string
}

type DayView struct {
	Date      time.Time
	DayName   string // "Lundi", "Mardi"...
	DateShort string // "14/09"
	IsToday   bool
	HasMeals  bool
	Cells     []MealCellView
}

type plannedMealsView struct {
	WeekStart time.Time
	Days      []DayView
}

func (app *Application) getPlannedMealsView(weekStart time.Time) (plannedMealsView, error) {
	meals, err := app.Meals.GetWeek(weekStart, weekStart.AddDate(0, 0, 6))
	if err != nil {
		return plannedMealsView{}, err
	}

	index := make(map[string]models.Meal, len(meals))
	for _, meal := range meals {
		key := meal.Date.Format("2006-01-02") + "|" + string(meal.MealSlot)
		index[key] = meal
	}

	today := time.Now().Format("2006-01-02")

	days := make([]DayView, 0, 7)
	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		day := DayView{
			Date:      date,
			DayName:   frenchDayFull(date),
			DateShort: date.Format("02/01"),
			IsToday:   date.Format("2006-01-02") == today,
		}
		for _, slot := range mealSlots {
			key := date.Format("2006-01-02") + "|" + string(slot)
			existing := index[key]
			if existing.Dish != "" {
				day.HasMeals = true
			}
			day.Cells = append(day.Cells, MealCellView{
				Date: date, MealSlot: slot, Dish: existing.Dish, Notes: existing.Notes,
			})
		}
		days = append(days, day)
	}

	return plannedMealsView{WeekStart: weekStart, Days: days}, nil
}

func mealDotClass(slot models.MealSlot) string {
	switch slot {
	case models.MealSlotBreakfast:
		return "breakfast"
	case models.MealSlotNoon:
		return "noon"
	default:
		return "dinner"
	}
}
