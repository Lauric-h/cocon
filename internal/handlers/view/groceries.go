package view

import (
	"cocon/internal/models"
	"cocon/internal/store"
)

type Groceries struct {
	Total     int
	Unchecked []models.GroceryItem
	Checked   []models.GroceryItem
}

func BuildGroceries(groceries *store.GroceryStore) (Groceries, error) {
	items, err := groceries.List()
	if err != nil {
		return Groceries{}, err
	}

	view := Groceries{Total: len(items)}
	for _, item := range items {
		if item.Checked {
			view.Checked = append(view.Checked, item)
		} else {
			view.Unchecked = append(view.Unchecked, item)
		}
	}
	return view, nil
}
