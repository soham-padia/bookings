package repository

import "github.com/solow-crypt/bookings/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testpassword string) (int, string, error)
	AddUser(first_name, last_name, email, password string) error
}
