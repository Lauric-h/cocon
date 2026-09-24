package view

import (
	"cocon/internal/calendar"
	"cocon/internal/helpers"
	"cocon/internal/models"
	"cocon/internal/store"
	"sort"
	"time"
)

const maxEventsPerDay = 1

type EventChip struct {
	Title     string
	TimeLabel string
}

type Event struct {
	Title     string
	DateLabel string
	IsPast    bool
}

type WeekDay struct {
	DayShort  string
	DayNum    string
	IsToday   bool
	MealDots  []string
	Events    []EventChip
	MoreCount int
}

type Dashboard struct {
	TodayLabel     string
	WeekRangeLabel string
	WeekDays       []WeekDay
	TodayMeals     []MealCell
	TodayEvents    []Event
	TasksRemaining int
	TasksPreview   []TaskItem
	GroceriesLeft  int
	Events         []Event
}

func BuildDashboard(meals *store.MealStore, groceries *store.GroceryStore, tasks *store.TaskStore, cal *calendar.Client) (Dashboard, error) {
	weekStart := helpers.StartOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 6)
	today := time.Now().Format("2006-01-02")

	mealsView, err := BuildPlannedMeals(meals, weekStart)
	if err != nil {
		return Dashboard{}, err
	}

	groceriesView, err := BuildGroceries(groceries)
	if err != nil {
		return Dashboard{}, err
	}

	tasksBlock, err := BuildDashboardTasks(tasks)
	if err != nil {
		return Dashboard{}, err
	}

	events, err := cal.EventsBetween(weekStart, weekEnd.AddDate(0, 0, 1).Add(-time.Second))
	if err != nil {
		return Dashboard{}, err
	}

	view := Dashboard{
		TodayLabel:     helpers.DateLong(time.Now()),
		WeekRangeLabel: weekStart.Format("2") + " – " + weekEnd.Format("2 Jan"),
		GroceriesLeft:  len(groceriesView.Unchecked),
		TasksRemaining: tasksBlock.TasksRemaining,
		TasksPreview:   tasksBlock.TasksPreview,
	}

	now := time.Now()
	var upcoming, past []models.Event
	for _, ev := range events {
		if ev.Start.Before(now) {
			past = append(past, ev)
		} else {
			upcoming = append(upcoming, ev)
		}
	}
	sort.Slice(upcoming, func(i, j int) bool { return upcoming[i].Start.Before(upcoming[j].Start) })
	sort.Slice(past, func(i, j int) bool { return past[i].Start.Before(past[j].Start) })

	eventsByDay := make(map[string][]EventChip)

	for _, ev := range upcoming {
		dateStr := ev.Start.Format("2006-01-02")

		chip := EventChip{
			Title:     ev.Title,
			TimeLabel: eventTimeLabel(ev),
		}

		eventsByDay[dateStr] = append(
			eventsByDay[dateStr],
			chip,
		)

		event := Event{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
		}

		view.Events = append(view.Events, event)

		if dateStr == today {
			view.TodayEvents = append(view.TodayEvents, event)
		}
	}

	for _, ev := range past {
		dateStr := ev.Start.Format("2006-01-02")

		chip := EventChip{
			Title:     ev.Title,
			TimeLabel: eventTimeLabel(ev),
		}

		eventsByDay[dateStr] = append(
			eventsByDay[dateStr],
			chip,
		)

		event := Event{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
			IsPast:    true,
		}

		view.Events = append(view.Events, event)

		if dateStr == today {
			view.TodayEvents = append(view.TodayEvents, event)
		}
	}

	for _, day := range mealsView.Days {
		dateStr := day.Date.Format("2006-01-02")
		dayEvents := eventsByDay[dateStr]

		strip := WeekDay{
			DayShort: helpers.DayShort(day.Date),
			DayNum:   day.Date.Format("2"),
			IsToday:  dateStr == today,
		}
		if len(dayEvents) > maxEventsPerDay {
			strip.Events = dayEvents[:maxEventsPerDay]
			strip.MoreCount = len(dayEvents) - maxEventsPerDay
		} else {
			strip.Events = dayEvents
		}

		for _, cell := range day.Cells {
			if cell.Dish == "" {
				continue
			}
			strip.MealDots = append(strip.MealDots, mealDotClass(cell.MealSlot))
			if strip.IsToday {
				view.TodayMeals = append(view.TodayMeals, cell)
			}
		}
		view.WeekDays = append(view.WeekDays, strip)
	}

	return view, nil
}

func eventTimeLabel(ev models.Event) string {
	if ev.AllDay {
		return "Toute la journée"
	}
	return ev.Start.Format("15h04")
}

func eventDateLabel(ev models.Event) string {
	return helpers.DayShort(ev.Start) + " " + ev.Start.Format("02/01") + " · " + eventTimeLabel(ev)
}
