package handlers

import (
	"net/http"

	"github.com/solow-crypt/bookings/pkg/config"
	"github.com/solow-crypt/bookings/pkg/models"
	"github.com/solow-crypt/bookings/pkg/render"
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

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Pc(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "pc.page.html", &models.TemplateData{})
}
func (m *Repository) Phone(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "phone.page.html", &models.TemplateData{})
}
func (m *Repository) Laptop(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "laptop.page.html", &models.TemplateData{})
}
func (m *Repository) Download(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "download.page.html", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{})
}
func (m *Repository) Docs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "docs.page.html", &models.TemplateData{})
}
func (m *Repository) Donate(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "donate.page.html", &models.TemplateData{})
}
