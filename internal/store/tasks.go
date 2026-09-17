package store

import (
	"cocon/internal/models"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type TaskStore struct {
	DB *sql.DB
}

// --- Task (la règle) ---

func (t *TaskStore) Create(task models.Task) (int, error) {
	res, err := t.DB.Exec(`
		INSERT INTO tasks (
			title, description, category, due_date,
			recurrence_frequency, recurrence_interval, recurrence_days_of_week,
			assigned_to, source
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Title, task.Description, string(task.TaskCategory), nullableDate(task.DueDate),
		recurrenceFrequency(task.Recurrence), recurrenceInterval(task.Recurrence), recurrenceDays(task.Recurrence),
		task.AssignedTo, task.Source,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (t *TaskStore) Update(task models.Task) error {
	_, err := t.DB.Exec(`
		UPDATE tasks SET
			title = ?, description = ?, category = ?, due_date = ?,
			recurrence_frequency = ?, recurrence_interval = ?, recurrence_days_of_week = ?,
			assigned_to = ?, source = ?
		WHERE id = ?
	`,
		task.Title, task.Description, string(task.TaskCategory), nullableDate(task.DueDate),
		recurrenceFrequency(task.Recurrence), recurrenceInterval(task.Recurrence), recurrenceDays(task.Recurrence),
		task.AssignedTo, task.Source, task.ID,
	)
	return err
}

func (t *TaskStore) Delete(id int) error {
	_, err := t.DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func (t *TaskStore) Get(id int) (models.Task, error) {
	row := t.DB.QueryRow(`
		SELECT id, title, description, category, due_date,
		       recurrence_frequency, recurrence_interval, recurrence_days_of_week,
		       assigned_to, source
		FROM tasks WHERE id = ?`, id)
	return scanTaskRow(row)
}

// ListBacklog renvoie les tâches "de fond" : ni date, ni récurrence.
func (t *TaskStore) ListBacklog() ([]models.Task, error) {
	return t.query(`
		SELECT id, title, description, category, due_date,
		       recurrence_frequency, recurrence_interval, recurrence_days_of_week,
		       assigned_to, source
		FROM tasks
		WHERE due_date IS NULL AND recurrence_frequency IS NULL
		ORDER BY id ASC`)
}

// ListDueBetween renvoie les tâches ponctuelles (non récurrentes) dont l'échéance tombe dans l'intervalle.
func (t *TaskStore) ListDueBetween(start, end time.Time) ([]models.Task, error) {
	return t.query(`
		SELECT id, title, description, category, due_date,
		       recurrence_frequency, recurrence_interval, recurrence_days_of_week,
		       assigned_to, source
		FROM tasks
		WHERE recurrence_frequency IS NULL AND due_date >= ? AND due_date <= ?
		ORDER BY due_date ASC`, start.Format("2006-01-02"), end.Format("2006-01-02"))
}

// ListRecurring renvoie toutes les tâches récurrentes.
func (t *TaskStore) ListRecurring() ([]models.Task, error) {
	return t.query(`
		SELECT id, title, description, category, due_date,
		       recurrence_frequency, recurrence_interval, recurrence_days_of_week,
		       assigned_to, source
		FROM tasks
		WHERE recurrence_frequency IS NOT NULL`)
}

func (t *TaskStore) query(query string, args ...any) ([]models.Task, error) {
	rows, err := t.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// scanTask lit une ligne depuis *sql.Rows (boucle sur plusieurs résultats).
func scanTask(rows *sql.Rows) (models.Task, error) {
	var task models.Task
	var category string
	var dueDate sql.NullString
	var recFreq sql.NullString
	var recInterval sql.NullInt64
	var recDays sql.NullString

	if err := rows.Scan(
		&task.ID, &task.Title, &task.Description, &category, &dueDate,
		&recFreq, &recInterval, &recDays,
		&task.AssignedTo, &task.Source,
	); err != nil {
		return task, err
	}

	task.TaskCategory = models.TaskCategory(category)
	task.DueDate = parseNullableDate(dueDate)
	task.Recurrence = buildRecurrence(recFreq, recInterval, recDays)

	return task, nil
}

// scanTaskRow lit une ligne depuis *sql.Row (résultat unique, ex: Get).
func scanTaskRow(row *sql.Row) (models.Task, error) {
	var task models.Task
	var category string
	var dueDate sql.NullString
	var recFreq sql.NullString
	var recInterval sql.NullInt64
	var recDays sql.NullString

	if err := row.Scan(
		&task.ID, &task.Title, &task.Description, &category, &dueDate,
		&recFreq, &recInterval, &recDays,
		&task.AssignedTo, &task.Source,
	); err != nil {
		return task, err
	}

	task.TaskCategory = models.TaskCategory(category)
	task.DueDate = parseNullableDate(dueDate)
	task.Recurrence = buildRecurrence(recFreq, recInterval, recDays)

	return task, nil
}

func buildRecurrence(freq sql.NullString, interval sql.NullInt64, days sql.NullString) *models.RecurrenceRule {
	if !freq.Valid {
		return nil
	}
	return &models.RecurrenceRule{
		Frequency:  freq.String,
		Interval:   int(interval.Int64),
		DaysOfWeek: decodeDaysOfWeek(days.String),
	}
}

// --- TaskOccurrence : porte la complétion pour TOUTE tâche (ponctuelle ou récurrente) ---

// EnsureOccurrence crée l'occurrence si elle n'existe pas encore pour cette date (idempotent).
func (t *TaskStore) EnsureOccurrence(taskID int, date time.Time) error {
	_, err := t.DB.Exec(`
		INSERT OR IGNORE INTO task_occurrences (task_id, date, assigned_to, done)
		VALUES (?, ?, '', 0)
	`, taskID, date.Format("2006-01-02"))
	return err
}

func (t *TaskStore) ListOccurrencesBetween(start, end time.Time) ([]models.TaskOccurrence, error) {
	rows, err := t.DB.Query(`
		SELECT id, task_id, date, assigned_to, done, completed_at
		FROM task_occurrences
		WHERE date >= ? AND date <= ?
		ORDER BY date ASC
	`, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []models.TaskOccurrence
	for rows.Next() {
		var occ models.TaskOccurrence
		var dateStr string
		var completedAt sql.NullString
		var done int
		if err := rows.Scan(&occ.ID, &occ.TaskID, &dateStr, &occ.AssignedTo, &done, &completedAt); err != nil {
			return nil, err
		}
		occ.Date, _ = time.Parse("2006-01-02", dateStr)
		occ.Done = done == 1
		occ.CompletedAt = parseNullableDate(completedAt)
		occurrences = append(occurrences, occ)
	}
	return occurrences, rows.Err()
}

func (t *TaskStore) MarkOccurrenceDone(id int) error {
	_, err := t.DB.Exec(`
		UPDATE task_occurrences SET done = 1, completed_at = ? WHERE id = ?
	`, time.Now().Format(time.RFC3339), id)
	return err
}

func (t *TaskStore) SetOccurrenceAssignee(id int, assignedTo string) error {
	_, err := t.DB.Exec(`UPDATE task_occurrences SET assigned_to = ? WHERE id = ?`, assignedTo, id)
	return err
}

// --- helpers ---

func nullableDate(d *time.Time) any {
	if d == nil {
		return nil
	}
	return d.Format("2006-01-02")
}

func parseNullableDate(s sql.NullString) *time.Time {
	if !s.Valid || s.String == "" {
		return nil
	}
	if parsed, err := time.Parse("2006-01-02", s.String); err == nil {
		return &parsed
	}
	if parsed, err := time.Parse(time.RFC3339, s.String); err == nil {
		return &parsed
	}
	return nil
}

func recurrenceFrequency(r *models.RecurrenceRule) any {
	if r == nil {
		return nil
	}
	return r.Frequency
}

func recurrenceInterval(r *models.RecurrenceRule) any {
	if r == nil {
		return nil
	}
	return r.Interval
}

func recurrenceDays(r *models.RecurrenceRule) any {
	if r == nil || len(r.DaysOfWeek) == 0 {
		return nil
	}
	parts := make([]string, len(r.DaysOfWeek))
	for i, d := range r.DaysOfWeek {
		parts[i] = strconv.Itoa(int(d))
	}
	return strings.Join(parts, ",")
}

func decodeDaysOfWeek(s string) []time.Weekday {
	if s == "" {
		return nil
	}
	var days []time.Weekday
	for _, p := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(p); err == nil {
			days = append(days, time.Weekday(n))
		}
	}
	return days
}
