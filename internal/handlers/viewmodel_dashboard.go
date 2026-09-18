package handlers

import (
	"cocon/internal/models"
	"sort"
	"time"
)

type WeekDayStripView struct {
	DayShort string // "LUN"
	DayNum   string // "14"
	IsToday  bool
	MealDots []string // classes CSS : "breakfast", "noon", "dinner"
	HasEvent bool
}

type dashboardTasksView struct {
	TasksRemaining int
	TasksPreview   []TaskItemView
}

type dashboardView struct {
	TodayLabel     string // "mercredi 17 septembre 2026"
	WeekRangeLabel string // "14 – 20 sept."
	WeekDays       []WeekDayStripView
	TodayMeals     []MealCellView // uniquement les cases remplies
	TasksRemaining int
	TasksPreview   []TaskItemView
	GroceriesLeft  int
	Events         []EventView
}

type EventView struct {
	Title     string
	DateLabel string
	IsPast    bool
}

func (app *Application) getDashboardView() (dashboardView, error) {
	weekStart := startOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 6)
	today := time.Now().Format("2006-01-02")

	mealsView, err := app.getPlannedMealsView(weekStart)
	if err != nil {
		return dashboardView{}, err
	}

	groceries, err := app.getGroceriesView()
	if err != nil {
		return dashboardView{}, err
	}

	tasksBlock, err := app.getDashboardTasksView()
	if err != nil {
		return dashboardView{}, err
	}

	events, err := app.Calendar.EventsBetween(weekStart, weekEnd.AddDate(0, 0, 1).Add(-time.Second))
	if err != nil {
		return dashboardView{}, err
	}

	view := dashboardView{
		TodayLabel:     frenchDateLong(time.Now()),
		WeekRangeLabel: weekStart.Format("2") + " – " + weekEnd.Format("2 Jan"),
		GroceriesLeft:  len(groceries.Unchecked),
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

	eventDaysWithEvent := make(map[string]bool)
	for _, ev := range upcoming {
		eventDaysWithEvent[ev.Start.Format("2006-01-02")] = true
		view.Events = append(view.Events, EventView{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
		})
	}
	for _, ev := range past {
		eventDaysWithEvent[ev.Start.Format("2006-01-02")] = true
		view.Events = append(view.Events, EventView{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
			IsPast:    true,
		})
	}

	for _, day := range mealsView.Days {
		dateStr := day.Date.Format("2006-01-02")
		strip := WeekDayStripView{
			DayShort: frenchDayShort(day.Date),
			DayNum:   day.Date.Format("2"),
			IsToday:  dateStr == today,
			HasEvent: eventDaysWithEvent[dateStr],
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

func eventDateLabel(ev models.Event) string {
	label := frenchDayShort(ev.Start) + " " + ev.Start.Format("02/01")
	if ev.AllDay {
		return label + " · Toute la journée"
	}
	return label + " · " + ev.Start.Format("15h04")
}

func (app *Application) getDashboardTasksView() (dashboardTasksView, error) {
	weekStart := startOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekStartStr, weekEndStr := weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")

	tasksView, err := app.getAllTasksView()
	if err != nil {
		return dashboardTasksView{}, err
	}

	view := dashboardTasksView{}
	for _, item := range tasksView.Planned {
		if item.DueDate == "" || item.DueDate < weekStartStr || item.DueDate > weekEndStr {
			continue
		}
		if !item.Done {
			view.TasksRemaining++
		}
		view.TasksPreview = append(view.TasksPreview, item)
	}
	if len(view.TasksPreview) > 5 {
		view.TasksPreview = view.TasksPreview[:5]
	}
	return view, nil
}
