package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/phone", handlers.Repo.Phone)
	mux.Get("/pc", handlers.Repo.Pc)
	mux.Get("/laptop", handlers.Repo.Laptop)
	mux.Get("/downloads", handlers.Repo.Download)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/docs", handlers.Repo.Docs)
	mux.Get("/donate", handlers.Repo.Donate)
	mux.Get("/login", handlers.Repo.Login)

	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/login-json", handlers.Repo.LoginJSON)

	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegistration)
	mux.Post("/register-json", handlers.Repo.RegisterJSON)
	mux.Get("/registration-summary", handlers.Repo.RegistrationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
