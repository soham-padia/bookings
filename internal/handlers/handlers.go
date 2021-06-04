package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/forms"
	"github.com/solow-crypt/bookings/internal/helpers"
	"github.com/solow-crypt/bookings/internal/models"
	"github.com/solow-crypt/bookings/internal/render"
)

//TemplateData holds data sent from handlers to templates

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Pc(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "pc.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Phone(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "phone.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Laptop(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "laptop.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Download(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "download.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Docs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "docs.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Donate(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "donate.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "log.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	w.Write([]byte(fmt.Sprintf("email: %s , password: %s", email, password)))
}

type jsonRequest struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) LoginJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonRequest{
		OK:      true,
		Message: "available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	var emptyRegistration models.Registration
	data := make(map[string]interface{})
	data["register"] = emptyRegistration

	render.RenderTemplate(w, r, "reg.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration := models.Registration{
		Firstname: r.Form.Get("first-name"),
		Lastname:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		Password:  r.Form.Get("password"),
		Password2: r.Form.Get("passwordre"),
	}

	form := forms.New(r.PostForm)

	form.Required("first-name", "last-name", "phone", "email", "password", "passwordre")
	form.MinLength("first-name", 3)
	form.IsEmail("email")
	form.IsPhoneNumber("phone")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["register"] = registration

		render.RenderTemplate(w, r, "reg.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "registration", registration)

	http.Redirect(w, r, "/registration-summary", http.StatusSeeOther)
}

func (m *Repository) RegisterJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonRequest{
		OK:      true,
		Message: "available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) RegistrationSummary(w http.ResponseWriter, r *http.Request) {

	registration, ok := m.App.Session.Get(r.Context(), "registration").(models.Registration)
	if !ok {
		m.App.ErrorLog.Println("Cant get error from session")
		m.App.Session.Put(r.Context(), "error", "Cant get registration summary")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "registration")

	data := make(map[string]interface{})
	data["registration"] = registration

	render.RenderTemplate(w, r, "reg.summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
