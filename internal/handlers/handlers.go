package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/forms"
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

	remoteIP := r.RemoteAddr
	//fmt.Println(remoteIP)
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Pc(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "pc.page.html", &models.TemplateData{})
}
func (m *Repository) Phone(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "phone.page.html", &models.TemplateData{})
}
func (m *Repository) Laptop(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "laptop.page.html", &models.TemplateData{})
}
func (m *Repository) Download(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "download.page.html", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}
func (m *Repository) Docs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "docs.page.html", &models.TemplateData{})
}
func (m *Repository) Donate(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "donate.page.html", &models.TemplateData{})
}
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "log.page.html", &models.TemplateData{})
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

	render.RenderTemplate(w, r, "reg.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
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
	form.MinLength("first-name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["register"] = registration

		render.RenderTemplate(w, r, "reg.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
}

func (m *Repository) RegisterJSON(w http.ResponseWriter, r *http.Request) {
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
