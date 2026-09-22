package view

import (
	"cocon/internal/helpers"
	"cocon/internal/models"
	"cocon/internal/store"
	"time"
)

var categoryLabels = map[models.TaskCategory]string{
	models.CategoryHouse:  "Maison",
	models.CategoryCar:    "Voiture",
	models.CategorySchool: "École",
	models.CategoryAdmin:  "Admin",
	models.CategoryFamily: "Famille",
	models.CategoryMeal:   "Repas",
}

type TaskItem struct {
	ID            int
	Title         string
	CategoryLabel string
	DueDate       string
	DueDateLabel  string
	Done          bool
	OccurrenceID  int
}

type AllTasks struct {
	Planned []TaskItem
	Backlog []TaskItem
}

func BuildAllTasks(tasks *store.TaskStore) (AllTasks, error) {
	planned, err := tasks.ListPlanned()
	if err != nil {
		return AllTasks{}, err
	}
	backlog, err := tasks.ListBacklog()
	if err != nil {
		return AllTasks{}, err
	}

	view := AllTasks{}
	for _, task := range planned {
		item, err := buildTaskItem(tasks, task)
		if err != nil {
			return AllTasks{}, err
		}
		view.Planned = append(view.Planned, item)
	}
	for _, task := range backlog {
		item, err := buildTaskItem(tasks, task)
		if err != nil {
			return AllTasks{}, err
		}
		view.Backlog = append(view.Backlog, item)
	}
	return view, nil
}

func buildTaskItem(tasks *store.TaskStore, task models.Task) (TaskItem, error) {
	item := TaskItem{
		ID:            task.ID,
		Title:         task.Title,
		CategoryLabel: categoryLabels[task.TaskCategory],
	}
	if task.DueDate != nil {
		item.DueDate = task.DueDate.Format("2006-01-02")
		item.DueDateLabel = helpers.DayFull(*task.DueDate) + " " + task.DueDate.Format("02/01")
	}

	occ, err := tasks.GetOccurrenceForTask(task.ID)
	if err != nil {
		return item, err
	}
	if occ != nil {
		item.OccurrenceID = occ.ID
		item.Done = occ.Done
	}
	return item, nil
}

type DashboardTasks struct {
	TasksRemaining int
	TasksPreview   []TaskItem
}

func BuildDashboardTasks(tasks *store.TaskStore) (DashboardTasks, error) {
	weekStart := helpers.StartOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekStartStr, weekEndStr := weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")

	all, err := BuildAllTasks(tasks)
	if err != nil {
		return DashboardTasks{}, err
	}

	view := DashboardTasks{}
	for _, item := range all.Planned {
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
