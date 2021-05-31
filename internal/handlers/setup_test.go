package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	// "github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/models"
	"github.com/solow-crypt/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.Registration{})

	//change this to true when in  production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //change it when going for release

	app.Session = session

	tc, err := render.CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template")
	}

	app.TemplateCache = tc

	//for developers use false
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/phone", Repo.Phone)
	mux.Get("/pc", Repo.Pc)
	mux.Get("/laptop", Repo.Laptop)
	mux.Get("/downloads", Repo.Download)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/docs", Repo.Docs)
	mux.Get("/donate", Repo.Donate)
	mux.Get("/login", Repo.Login)

	mux.Post("/login", Repo.PostLogin)
	mux.Get("/login-json", Repo.LoginJSON)

	mux.Get("/register", Repo.Register)
	mux.Post("/register", Repo.PostRegistration)
	mux.Post("/register-json", Repo.RegisterJSON)
	mux.Get("/registration-summary", Repo.RegistrationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

//adds CSRF protection to all POST requests
func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
