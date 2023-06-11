package model

import (
	"barlio/internal/data"
	"barlio/internal/validator"
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                        int
	Firstname                 data.String
	Lastname                  data.String
	Username                  data.String
	Email                     data.String
	Password                  data.String
	IsVerified                bool
	JoinedAt                  time.Time
	Birthdate                 time.Time
	PreferedArticleCategories ListArticleCategorie
	Bio                       data.String
	ProfilPicture             data.String
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*&u.Password = data.String(hashedPassword)
	return nil
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(u *User) error {
	const statement = `INSERT INTO users(
							username, password, email
						) VALUES ($1, $2, $3) 
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return m.DB.QueryRowContext(ctx, statement).Scan(&u.ID)
}

func (m *UserModel) Get(u User) (*User, error) {
	var user User
	const statement = `SELECT 
							firstname, lastname, username,  password, email, joined_at,
							isverified
						FROM users
						WHERE 
							(id = $1 || $1 = 0)
							AND (username=$2 || $2 = '')
							AND (email=$3 || $3 = '')`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, statement, u.ID, u.Username, u.Email).Scan(&u.Firstname, &u.Lastname, &u.Username,
		&u.Password, &u.Email, &u.JoinedAt, &u.IsVerified)

	return &user, err
}

func (m *UserModel) Delete(id int) error {
	const statement = `UPDATE users
						SET 
							deleted_request_date=$1
						WHERE 
							id=$1
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return m.DB.QueryRowContext(ctx, statement, id).Scan(&id)
}

func (m *UserModel) GetAll(id int) ([]User, error) {
	var users []User
	const statement = `SELECT 
							username, firstname, lastname, profil_picture
						FROM users
						WHERE
							isverified=true`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User

		err = rows.Scan(&u.Username, &u.Firstname, &u.Lastname, &u.ProfilPicture)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (m *UserModel) ValidateUser(u *User, validator *validator.Validator) error {
	user, err := m.Get(User{Username: data.String(u.Username)})
	if err != nil {
		return err
	}
	validator.Check(user == nil, "username", "username already taken")

	user, err = m.Get(User{Email: data.String(u.Email)})
	if err != nil {
		return err
	}
	validator.Check(user == nil, "email", "email already in use")
	return nil
}
