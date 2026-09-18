package handlers

import "time"

type WeekDayStripView struct {
	DayShort string // "LUN"
	DayNum   string // "14"
	IsToday  bool
	MealDots []string // classes CSS : "breakfast", "noon", "dinner"
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

	view := dashboardView{
		TodayLabel:     frenchDateLong(time.Now()),
		WeekRangeLabel: weekStart.Format("2") + " – " + weekEnd.Format("2 Jan"),
		GroceriesLeft:  len(groceries.Unchecked),
	}

	for _, day := range mealsView.Days {
		strip := WeekDayStripView{
			DayShort: frenchDayShort(day.Date),
			DayNum:   day.Date.Format("2"),
			IsToday:  day.Date.Format("2006-01-02") == today,
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

	tasksBlock, err := app.getDashboardTasksView()
	if err != nil {
		return dashboardView{}, err
	}
	view.TasksRemaining = tasksBlock.TasksRemaining
	view.TasksPreview = tasksBlock.TasksPreview

	return view, nil
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
