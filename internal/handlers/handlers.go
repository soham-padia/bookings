package handlers

import (
	"log"
	"net/http"

	"github.com/solow-crypt/bookings/internal/config"
	"github.com/solow-crypt/bookings/internal/driver"
	"github.com/solow-crypt/bookings/internal/forms"
	"github.com/solow-crypt/bookings/internal/models"
	"github.com/solow-crypt/bookings/internal/render"
	"github.com/solow-crypt/bookings/internal/repository"
	"github.com/solow-crypt/bookings/internal/repository/dbrepo"
)

//TemplateData holds data sent from handlers to templates

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Pc(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "pc.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Phone(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "phone.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Laptop(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "laptop.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Download(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "download.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Docs(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "docs.page.tmpl", &models.TemplateData{})
}
func (m *Repository) Donate(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "donate.page.tmpl", &models.TemplateData{})
}

// type jsonRequest struct {
// 	OK      bool   `json:"ok"`
// 	Message string `json:"message"`
// }

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "log.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//handles the login of user
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "log.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)

	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository) Registration(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "reg.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	passwordre := r.Form.Get("passwordre")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	form.IsSame("password", "passwordre")

	if !form.Valid() {
		render.Template(w, r, "reg.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	err = m.DB.AddUser(first_name, last_name, email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid registration credentials")
		http.Redirect(w, r, "/user/register", http.StatusSeeOther)
		return
	}
	// m.App.Session.Put(r.Context(), "user_id", id)

	m.App.Session.Put(r.Context(), "flash", "registered successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminNewUsers(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-new-users.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminAllUsers(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "admin-all-users.page.tmpl", &models.TemplateData{})
}
func (m *Repository) AdminDonationInfo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-donations-info.page.tmpl", &models.TemplateData{})
}
