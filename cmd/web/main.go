package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	calendarpkg "cocon/internal/calendar"
	"cocon/internal/handlers"
	"cocon/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Println("pas de fichier .env trouvé, on continue avec l'environnement système")
	}

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

	calendarClient := calendarpkg.NewClient(os.Getenv("GOOGLE_CALENDAR_ICS_URL"))
	log.Printf("%s\n", os.Getenv("GOOGLE_CALENDAR_ICS_URL"))

	if events, err := calendarClient.EventsBetween(time.Now(), time.Now().AddDate(0, 0, 7)); err != nil {
		log.Println("erreur calendrier:", err)
	} else {
		log.Printf("%d événements trouvés cette semaine: %+v", len(events), events)
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
		Calendar:  calendarClient,
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
