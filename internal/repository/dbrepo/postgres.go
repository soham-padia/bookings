package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/solow-crypt/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//returns a user by id
func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at 
				from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

//updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update user set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

//authenticates a user
func (m *postgresDBRepo) Authenticate(email, testpassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testpassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

func (m *postgresDBRepo) AddUser(first_name, last_name, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	created_at := time.Now()
	updated_at := time.Now()

	checkforEmail := m.DB.QueryRowContext(ctx, "select email from users where email = $1", email)

	var email2 string

	err := checkforEmail.Scan(&email2)
	if err != nil {
		if email2 == email {
			log.Println(err)
			return errors.New("email already exists")
		}
	}

	query := `insert into users (first_name,last_name,email,password,access_level,created_at,updated_at) values ($1, $2, $3, $4, $5, $6);`
	row := m.DB.QueryRowContext(ctx, query, first_name, last_name, email, hashedPassword, 1, created_at, updated_at)
	err = row.Scan()
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) InsertUser(u models.Registration) (error, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	checkforEmail := `select id from users where email='$1'`

	id, err := m.DB.ExecContext(ctx, checkforEmail, u.Email)
	if err != nil {

		if !(id != nil) {
			//TODO
			log.Println("------------------------------------------------")
			return err, 1
		}

		return err, 0
	}

	stmt := `insert into users (first_name,last_name,email,password,access_level, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 12)

	_, err = m.DB.ExecContext(ctx, stmt,
		u.Firstname,
		u.Lastname,
		u.Email,
		hashedPassword,
		1,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err, 0
	}

	return nil, 0
}
