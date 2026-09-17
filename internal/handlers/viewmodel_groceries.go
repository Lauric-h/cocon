package handlers

import "cocon/internal/models"

type groceriesView struct {
	Total     int
	Unchecked []models.GroceryItem
	Checked   []models.GroceryItem
}

func buildGroceriesView(items []models.GroceryItem) groceriesView {
	view := groceriesView{Total: len(items)}
	for _, item := range items {
		if item.Checked {
			view.Checked = append(view.Checked, item)
		} else {
			view.Unchecked = append(view.Unchecked, item)
		}
	}
	return view
}

func (app *Application) getGroceriesView() (groceriesView, error) {
	items, err := app.Groceries.List()
	if err != nil {
		return groceriesView{}, err
	}
	return buildGroceriesView(items), nil
}
