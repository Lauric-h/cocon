package store

import (
	"cocon/internal/models"
	"database/sql"
	"time"
)

type MealStore struct {
	DB *sql.DB
}

func (m *MealStore) GetWeek(start, end time.Time) ([]models.Meal, error) {
	rows, err := m.DB.Query(`
		SELECT id, date, meal_slot, dish, notes
		FROM meals WHERE date >= ? AND date <= ?
		ORDER BY date ASC, 
		    CASE meal_slot
			WHEN 'petit-déjeuner' THEN 1
			WHEN 'midi' THEN 2
			WHEN 'dîner' THEN 3
		END ASC
	`, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		var dateStr, mealSlot string
		if err := rows.Scan(&meal.ID, &dateStr, &mealSlot, &meal.Dish, &meal.Notes); err != nil {
			return nil, err
		}
		meal.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		meal.MealSlot = models.MealSlot(mealSlot)
		meals = append(meals, meal)
	}
	return meals, rows.Err()
}

func (m *MealStore) Upsert(date time.Time, mealSlot models.MealSlot, dish, notes string) error {
	_, err := m.DB.Exec(`
		INSERT INTO meals (date, meal_slot, dish, notes)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(date, meal_slot) DO UPDATE SET dish = excluded.dish, notes = excluded.notes
	`, date.Format("2006-01-02"), string(mealSlot), dish, notes)
	return err
}

func (m *MealStore) Delete(id int) error {
	_, err := m.DB.Exec(`DELETE FROM meals WHERE id = ?`, id)
	return err
}
