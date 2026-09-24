package view

import (
	"cocon/internal/calendar"
	"cocon/internal/helpers"
	"cocon/internal/models"
	"cocon/internal/store"
	"time"
)

type WeeklySummaryEvent struct {
	Title     string
	TimeLabel string // "17h00" ou "Toute la journée"
}

type WeeklySummaryDay struct {
	DayName   string
	DateLabel string
	Events    []WeeklySummaryEvent
	Breakfast string
	Lunch     string
	Dinner    string
}

func BuildWeeklySummaryDays(meals *store.MealStore, cal *calendar.Client, weekStart time.Time) ([]WeeklySummaryDay, error) {
	weekEnd := weekStart.AddDate(0, 0, 6)

	mealRows, err := meals.GetWeek(weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	mealIndex := make(map[string]models.Meal, len(mealRows))
	for _, m := range mealRows {
		mealIndex[m.Date.Format("2006-01-02")+"|"+string(m.MealSlot)] = m
	}

	events, err := cal.EventsBetween(weekStart, weekEnd.AddDate(0, 0, 1).Add(-time.Second))
	if err != nil {
		return nil, err
	}
	eventsByDay := make(map[string][]WeeklySummaryEvent)
	for _, ev := range events {
		dateStr := ev.Start.Format("2006-01-02")
		eventsByDay[dateStr] = append(eventsByDay[dateStr], WeeklySummaryEvent{
			Title:     ev.Title,
			TimeLabel: eventTimeLabel(ev),
		})
	}

	days := make([]WeeklySummaryDay, 0, 7)
	for i := 0; i < 7; i++ {
		date := weekStart.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		days = append(days, WeeklySummaryDay{
			DayName:   helpers.DayFull(date),
			DateLabel: helpers.DayMonth(date),
			Events:    eventsByDay[dateStr],
			Breakfast: mealIndex[dateStr+"|"+string(models.MealSlotBreakfast)].Dish,
			Lunch:     mealIndex[dateStr+"|"+string(models.MealSlotNoon)].Dish,
			Dinner:    mealIndex[dateStr+"|"+string(models.MealSlotDinner)].Dish,
		})
	}
	return days, nil
}
