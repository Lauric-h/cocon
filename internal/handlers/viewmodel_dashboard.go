package handlers

import (
	"cocon/internal/models"
	"sort"
	"time"
)

type EventChipView struct {
	Title     string
	TimeLabel string
}

type WeekDayStripView struct {
	DayShort  string
	DayNum    string
	IsToday   bool
	MealDots  []string
	Events    []EventChipView
	MoreCount int
}

const maxEventsPerDay = 1

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

	eventsByDay := make(map[string][]EventChipView)
	for _, ev := range upcoming {
		dateStr := ev.Start.Format("2006-01-02")
		eventsByDay[dateStr] = append(eventsByDay[dateStr], EventChipView{
			Title:     ev.Title,
			TimeLabel: eventTimeLabel(ev),
		})
		view.Events = append(view.Events, EventView{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
		})
	}
	for _, ev := range past {
		dateStr := ev.Start.Format("2006-01-02")
		eventsByDay[dateStr] = append(eventsByDay[dateStr], EventChipView{
			Title:     ev.Title,
			TimeLabel: eventTimeLabel(ev),
		})
		view.Events = append(view.Events, EventView{
			Title:     ev.Title,
			DateLabel: eventDateLabel(ev),
			IsPast:    true,
		})
	}

	for _, day := range mealsView.Days {
		dateStr := day.Date.Format("2006-01-02")
		dayEvents := eventsByDay[dateStr]

		strip := WeekDayStripView{
			DayShort: frenchDayShort(day.Date),
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
	return frenchDayShort(ev.Start) + " " + ev.Start.Format("02/01") + " · " + eventTimeLabel(ev)
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
