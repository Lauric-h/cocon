package handlers

import (
	"cocon/internal/models"
	"time"
)

var taskCategoryLabels = map[models.TaskCategory]string{
	models.CategoryHouse:  "Maison",
	models.CategoryCar:    "Voiture",
	models.CategorySchool: "École",
	models.CategoryAdmin:  "Admin",
	models.CategoryFamily: "Famille",
	models.CategoryMeal:   "Repas",
}

type TaskItemView struct {
	ID            int
	Title         string
	CategoryLabel string
	Done          bool
	OccurrenceID  int // 0 si pas d'occurrence (ne devrait pas arriver pour une tâche avec échéance)
}

type weekTasksView struct {
	Items []TaskItemView
}

func (app *Application) getWeekTasksView(weekStart time.Time) (weekTasksView, error) {
	tasks, err := app.Tasks.ListDueBetween(weekStart, weekStart.AddDate(0, 0, 6))
	if err != nil {
		return weekTasksView{}, err
	}

	var items []TaskItemView
	for _, task := range tasks {
		occ, err := app.Tasks.GetOccurrenceForTask(task.ID)
		if err != nil {
			return weekTasksView{}, err
		}
		item := TaskItemView{
			ID:            task.ID,
			Title:         task.Title,
			CategoryLabel: taskCategoryLabels[task.TaskCategory],
		}
		if occ != nil {
			item.OccurrenceID = occ.ID
			item.Done = occ.Done
		}
		items = append(items, item)
	}

	return weekTasksView{Items: items}, nil
}
