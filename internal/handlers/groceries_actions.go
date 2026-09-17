package handlers

import (
	"net/http"
	"strconv"
)

func (app *Application) createGrocery(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulaire invalide", http.StatusBadRequest)
		return
	}
	name := r.PostForm.Get("name")
	if name == "" {
		http.Error(w, "nom requis", http.StatusBadRequest)
		return
	}

	quantity, err := parseOptionalQuantity(r.PostForm.Get("quantity"))
	if err != nil {
		http.Error(w, "quantité invalide", http.StatusBadRequest)
		return
	}

	if err := app.Groceries.Create(name, quantity); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderGroceryList(w)
}

func (app *Application) toggleGrocery(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if err := app.Groceries.Toggle(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderGroceryList(w)
}

func (app *Application) deleteGrocery(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	if err := app.Groceries.Delete(id); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderGroceryList(w)
}

func (app *Application) clearCheckedGroceries(w http.ResponseWriter, r *http.Request) {
	if err := app.Groceries.ClearChecked(); err != nil {
		app.serverError(w, err)
		return
	}
	app.renderGroceryList(w)
}

func (app *Application) renderGroceryList(w http.ResponseWriter) {
	view, err := app.getGroceriesView()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "grocery-list.tmpl", view)
}

func parseOptionalQuantity(raw string) (*int, error) {
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
