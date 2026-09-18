package handlers

import (
	"cocon/internal/models"
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
	DueDate       string // "2006-01-02", "" si backlog — pour filtrage/tri
	DueDateLabel  string // "Mardi 15/09" — pour affichage
	Done          bool
	OccurrenceID  int
}

type allTasksView struct {
	Planned []TaskItemView // ont une échéance
	Backlog []TaskItemView // pas de date, "de fond"
}

func (app *Application) getAllTasksView() (allTasksView, error) {
	planned, err := app.Tasks.ListPlanned()
	if err != nil {
		return allTasksView{}, err
	}
	backlog, err := app.Tasks.ListBacklog()
	if err != nil {
		return allTasksView{}, err
	}

	view := allTasksView{}

	for _, task := range planned {
		item, err := app.buildTaskItemView(task)
		if err != nil {
			return allTasksView{}, err
		}
		view.Planned = append(view.Planned, item)
	}
	for _, task := range backlog {
		item, err := app.buildTaskItemView(task)
		if err != nil {
			return allTasksView{}, err
		}
		view.Backlog = append(view.Backlog, item)
	}

	return view, nil
}

func (app *Application) buildTaskItemView(task models.Task) (TaskItemView, error) {
	item := TaskItemView{
		ID:            task.ID,
		Title:         task.Title,
		CategoryLabel: taskCategoryLabels[task.TaskCategory],
	}
	if task.DueDate != nil {
		item.DueDate = task.DueDate.Format("2006-01-02")
		item.DueDateLabel = frenchDayFull(*task.DueDate) + " " + task.DueDate.Format("02/01")
	}

	occ, err := app.Tasks.GetOccurrenceForTask(task.ID)
	if err != nil {
		return item, err
	}
	if occ != nil {
		item.OccurrenceID = occ.ID
		item.Done = occ.Done
	}

	return item, nil
}
