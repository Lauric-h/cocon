package models

import "time"

type MealSlot string

const (
	MealSlotBreakfast MealSlot = "petit-déjeuner"
	MealSlotNoon      MealSlot = "midi"
	MealSlotDinner    MealSlot = "dîner"
)

type Meal struct {
	ID       int
	Date     time.Time
	MealSlot MealSlot
	Dish     string
	Notes    string
}
