package store

import (
	"cocon/internal/models"
	"database/sql"
	"time"
)

type GroceryStore struct {
	DB *sql.DB
}

func (g *GroceryStore) List() ([]models.GroceryItem, error) {
	rows, err := g.DB.Query(`
		SELECT id, name, quantity, checked, created_at
		FROM groceries
		ORDER BY checked ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.GroceryItem
	for rows.Next() {
		var it models.GroceryItem
		var quantity sql.NullInt64
		var createdStr string
		var checked int
		if err := rows.Scan(&it.ID, &it.Name, &quantity, &checked, &createdStr); err != nil {
			return nil, err
		}
		if quantity.Valid {
			it.Quantity = new(int(quantity.Int64))
		}
		it.Checked = checked == 1
		it.CreatedAt, err = time.Parse(time.RFC3339, createdStr)
		if err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (g *GroceryStore) Create(name string, quantity *int) error {
	_, err := g.DB.Exec(`
		INSERT INTO groceries (name, quantity, checked, created_at)
		VALUES (?, ?, 0, ?)
	`, name, nullableInt(quantity), time.Now().Format(time.RFC3339))
	return err
}

func (g *GroceryStore) Toggle(id int) error {
	_, err := g.DB.Exec(`UPDATE groceries SET checked = 1 - checked WHERE id = ?`, id)
	return err
}

func (g *GroceryStore) Delete(id int) error {
	_, err := g.DB.Exec(`DELETE FROM groceries WHERE id = ?`, id)
	return err
}

func (g *GroceryStore) ClearChecked() error {
	_, err := g.DB.Exec(`DELETE FROM groceries WHERE checked = 1`)
	return err
}

func nullableInt(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}
