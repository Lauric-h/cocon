package store

import (
	"cocon/internal/models"
	"database/sql"
	"time"
)

type TaskStore struct {
	DB *sql.DB
}

// Create insère la tâche, et si elle a une échéance, crée immédiatement son unique occurrence.
func (t *TaskStore) Create(task models.Task) (int, error) {
	res, err := t.DB.Exec(`
		INSERT INTO tasks (title, description, category, due_date, assigned_to, source)
		VALUES (?, ?, ?, ?, ?, ?)
	`, task.Title, task.Description, string(task.TaskCategory), nullableDate(task.DueDate), task.AssignedTo, task.Source)
	if err != nil {
		return 0, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	id := int(id64)

	if task.DueDate != nil {
		if err := t.EnsureOccurrence(id, *task.DueDate); err != nil {
			return id, err
		}
	}
	return id, nil
}

func (t *TaskStore) Delete(id int) error {
	_, err := t.DB.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func (t *TaskStore) Get(id int) (models.Task, error) {
	row := t.DB.QueryRow(`
		SELECT id, title, description, category, due_date, assigned_to, source
		FROM tasks WHERE id = ?`, id)
	return scanTaskRow(row)
}

// ListBacklog renvoie les tâches sans échéance.
func (t *TaskStore) ListBacklog() ([]models.Task, error) {
	return t.query(`
		SELECT id, title, description, category, due_date, assigned_to, source
		FROM tasks WHERE due_date IS NULL ORDER BY id ASC`)
}

// ListDueBetween renvoie les tâches dont l'échéance tombe dans l'intervalle.
func (t *TaskStore) ListDueBetween(start, end time.Time) ([]models.Task, error) {
	return t.query(`
		SELECT id, title, description, category, due_date, assigned_to, source
		FROM tasks
		WHERE due_date >= ? AND due_date <= ?
		ORDER BY due_date ASC`, start.Format("2006-01-02"), end.Format("2006-01-02"))
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

func scanTask(rows *sql.Rows) (models.Task, error) {
	var task models.Task
	var category string
	var dueDate sql.NullString
	if err := rows.Scan(&task.ID, &task.Title, &task.Description, &category, &dueDate, &task.AssignedTo, &task.Source); err != nil {
		return task, err
	}
	task.TaskCategory = models.TaskCategory(category)
	task.DueDate = parseNullableDate(dueDate)
	return task, nil
}

func scanTaskRow(row *sql.Row) (models.Task, error) {
	var task models.Task
	var category string
	var dueDate sql.NullString
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &category, &dueDate, &task.AssignedTo, &task.Source); err != nil {
		return task, err
	}
	task.TaskCategory = models.TaskCategory(category)
	task.DueDate = parseNullableDate(dueDate)
	return task, nil
}

// --- TaskOccurrence : porte la complétion (une seule occurrence par tâche ponctuelle) ---

func (t *TaskStore) EnsureOccurrence(taskID int, date time.Time) error {
	_, err := t.DB.Exec(`
		INSERT OR IGNORE INTO task_occurrences (task_id, date, assigned_to, done)
		VALUES (?, ?, '', 0)
	`, taskID, date.Format("2006-01-02"))
	return err
}

// GetOccurrenceForTask renvoie l'occurrence d'une tâche (une tâche ponctuelle n'en a qu'une).
func (t *TaskStore) GetOccurrenceForTask(taskID int) (*models.TaskOccurrence, error) {
	row := t.DB.QueryRow(`
		SELECT id, task_id, date, assigned_to, done, completed_at
		FROM task_occurrences WHERE task_id = ? LIMIT 1
	`, taskID)

	var occ models.TaskOccurrence
	var dateStr string
	var completedAt sql.NullString
	var done int
	if err := row.Scan(&occ.ID, &occ.TaskID, &dateStr, &occ.AssignedTo, &done, &completedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	occ.Date, _ = time.Parse("2006-01-02", dateStr)
	occ.Done = done == 1
	occ.CompletedAt = parseNullableDate(completedAt)
	return &occ, nil
}

func (t *TaskStore) ToggleOccurrenceDone(id int) error {
	_, err := t.DB.Exec(`
		UPDATE task_occurrences
		SET done = 1 - done, completed_at = CASE WHEN done = 0 THEN ? ELSE NULL END
		WHERE id = ?
	`, time.Now().Format(time.RFC3339), id)
	return err
}

// --- helpers (partagés avec store/meals.go — déjà présents si nullableDate/parseNullableDate existent) ---

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
