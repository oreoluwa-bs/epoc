package web

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oreoluwa-bs/epoc/web/handlers"
)

type App struct {
	config Config
	router http.Handler
	db     *sql.DB
}

func New(cfg Config) *App {

	app := &App{
		config: cfg,
	}
	app.loadDB()
	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.PORT),
		Handler: a.router,
	}

	_, err := a.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer func() {
		if err := a.db.Close(); err != nil {
			fmt.Println("failed to close database", err)
		}
	}()

	ch := make(chan error, 1)

	go func() {
		fmt.Println("Listening on port ", server.Addr)
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

	// return nil
}

func (a *App) loadDB() {
	db, err := sql.Open("sqlite3", a.config.DatabaseURL)

	if err != nil {
		log.Fatal("failed to initialise database: ", err)
	}

	// run migrations
	const eventMigration = `
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER NOT NULL PRIMARY KEY,
			name VARCHAR NOT NULL,
			description TEXT,
			starts_at DATETIME NOT NULL,
			ends_at DATETIME NOT NULL
		)
	`

	if _, err := db.Exec(eventMigration); err != nil {
		log.Panic(err)
	}

	a.db = db
}

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/events", a.loadEventsRouters)

	a.router = router
}

func (a *App) loadEventsRouters(r chi.Router) {
	handler := &handlers.EventHandler{
		DB: a.db,
	}

	r.Post("/", handler.Create)
	r.Get("/", handler.List)
	r.Get("/{id}", handler.GetById)
	r.Put("/{id}", handler.UpdateById)
	r.Delete("/{id}", handler.DeleteById)
}
