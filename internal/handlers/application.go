package handlers

import (
	"bytes"
	"cocon/internal/store"
	"html/template"
	"log"
	"net/http"
)

type Application struct {
	Meals     *store.MealStore
	Groceries *store.GroceryStore
	Templates *template.Template
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))

	mux.HandleFunc("GET /planned-meals", app.plannedMeals)
	mux.HandleFunc("POST /meals", app.upsertMeal)
	mux.HandleFunc("POST /meals/{id}/delete", app.deleteMeal)

	mux.HandleFunc("POST /groceries", app.createGrocery)
	mux.HandleFunc("POST /groceries/{id}/toggle", app.toggleGrocery)
	mux.HandleFunc("POST /groceries/{id}/delete", app.deleteGrocery)
	mux.HandleFunc("POST /groceries/clear-checked", app.clearCheckedGroceries)

	return mux
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "Erreur serveur", http.StatusInternalServerError)
}

func (app *Application) render(w http.ResponseWriter, status int, name string, data any) {
	var body bytes.Buffer
	if err := app.Templates.ExecuteTemplate(&body, name, data); err != nil {
		log.Println(err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(body.Bytes())
}
