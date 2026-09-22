package view

import (
	"cocon/internal/helpers"
	"cocon/internal/models"
	"cocon/internal/store"
	"time"
)

var mealSlots = []models.MealSlot{
	models.MealSlotBreakfast,
	models.MealSlotNoon,
	models.MealSlotDinner,
}

type MealCell struct {
	Date     time.Time
	MealSlot models.MealSlot
	Dish     string
	Notes    string
}

type Day struct {
	Date      time.Time
	DayName   string
	DateShort string
	IsToday   bool
	HasMeals  bool
	Cells     []MealCell
}

type PlannedMeals struct {
	WeekStart time.Time
	Days      []Day
}

func BuildPlannedMeals(meals *store.MealStore, weekStart time.Time) (PlannedMeals, error) {
	weekMeals, err := meals.GetWeek(weekStart, weekStart.AddDate(0, 0, 6))
	if err != nil {
		return PlannedMeals{}, err
	}

	index := make(map[string]models.Meal, len(weekMeals))
	for _, meal := range weekMeals {
		key := meal.Date.Format("2006-01-02") + "|" + string(meal.MealSlot)
		index[key] = meal
	}

	today := time.Now().Format("2006-01-02")

	days := make([]Day, 0, 7)
	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		day := Day{
			Date:      date,
			DayName:   helpers.DayFull(date),
			DateShort: date.Format("02/01"),
			IsToday:   date.Format("2006-01-02") == today,
		}
		for _, slot := range mealSlots {
			key := date.Format("2006-01-02") + "|" + string(slot)
			existing := index[key]
			if existing.Dish != "" {
				day.HasMeals = true
			}
			day.Cells = append(day.Cells, MealCell{
				Date: date, MealSlot: slot, Dish: existing.Dish, Notes: existing.Notes,
			})
		}
		days = append(days, day)
	}

	return PlannedMeals{WeekStart: weekStart, Days: days}, nil
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
