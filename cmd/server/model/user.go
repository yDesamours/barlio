package model

import (
	"barlio/internal/data"
	"database/sql"
	"time"
)

type User struct {
	ID                        int                  `json:"id,omitempty"`
	Firstname                 data.String          `json:"firstname,omitempty"`
	Lastname                  data.String          `json:"lastname,omitempty"`
	Username                  data.String          `json:"username,omitempty"`
	Email                     data.String          `json:"email,omitempty"`
	Password                  data.String          `json:"-"`
	IsVerified                bool                 `json:"-"`
	JoinedAt                  time.Time            `json:"joinedAt"`
	Birthdate                 time.Time            `json:"birthdate"`
	PreferedArticleCategories ListArticleCategorie `json:"preferedArticleCategories"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(u *User) error {
	const statement = `INSERT INTO users(
							username, password, email
						) VALUES ($1, $2, $3) 
						RETURNING id`

	return m.DB.QueryRow(statement).Scan(&u.ID)
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

	err := m.DB.QueryRow(statement, u.ID, u.Username, u.Email).Scan(&u.Firstname, &u.Lastname, &u.Username,
		&u.Password, &u.Email, &u.JoinedAt, &u.IsVerified)

	return &user, err
}

func (m *UserModel) Delete(id int) error {
	return nil
}

func (m *UserModel) Validate(id int) error {
	return nil
}
