package models

import "time"

type TaskCategory string

const (
	CategoryHouse  TaskCategory = "house"
	CategoryCar    TaskCategory = "car"
	CategorySchool TaskCategory = "school"
	CategoryAdmin  TaskCategory = "admin"
	CategoryFamily TaskCategory = "family"
	CategoryMeal   TaskCategory = "meal"
)

type RecurrenceRule struct {
	Frequency  string
	Interval   int
	DaysOfWeek []time.Weekday
}

type Task struct {
	ID           int
	Title        string
	Description  string
	TaskCategory TaskCategory
	DueDate      *time.Time
	Recurrence   *RecurrenceRule
	AssignedTo   string
	Source       string
}

func (t Task) IsBacklog() bool {
	return t.DueDate == nil && t.Recurrence == nil
}

type TaskOccurrence struct {
	ID          int
	TaskID      int
	Date        time.Time
	AssignedTo  string
	Done        bool
	CompletedAt *time.Time
}
