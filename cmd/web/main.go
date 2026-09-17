package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cocon/internal/handlers"
	"cocon/internal/store"
)

func main() {
	ctx := context.Background()
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/cocon.db"
	}
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatal(err)
		}
	}

	db, err := store.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	templates, err := template.ParseGlob("ui/html/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := templates.ParseGlob("ui/html/partials/*.tmpl"); err != nil {
		log.Fatal(err)
	}

	app := &handlers.Application{
		Meals:     &store.MealStore{DB: db},
		Groceries: &store.GroceryStore{DB: db},
		Templates: templates,
		Tasks:     &store.TaskStore{DB: db},
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	log.Printf("listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}
