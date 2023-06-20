package model

import (
	"barlio/internal/types"
	"context"
	"database/sql"
	"net/url"
	"time"
)

type User struct {
	ID                        int
	Firstname                 types.String
	Lastname                  types.String
	Username                  types.String
	Email                     types.String
	Password                  types.String
	IsVerified                bool
	JoinedAt                  time.Time
	Birthdate                 time.Time
	PreferedArticleCategories ListArticleCategorie
	Bio                       types.String
	ProfilPicture             types.String
}

func NullUser() User {
	return User{}
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) NewUser(form url.Values) *User {
	user := &User{
		Username: types.String(form.Get("username")),
		Email:    types.String(form.Get("email")),
		Password: types.String(form.Get("password")),
	}
	return user
}

func (m *UserModel) Insert(u *User) error {
	const statement = `INSERT INTO users(
							username, password, email, joined_at, isverified
						) VALUES ($1, $2, $3, $4, $5) 
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return m.DB.QueryRowContext(ctx, statement, u.Username, u.Password, u.Email, u.JoinedAt,
		u.IsVerified).Scan(&u.ID)
}

func (m *UserModel) Get(user User) (*User, error) {
	var u User
	const statement = `SELECT 
							id, firstname, lastname, username,  password, email, joined_at,
							isverified
						FROM users
						WHERE 
							(id = $1 OR $1 = 0)
							AND (trim(username)=$2 OR $2 is null)
							AND (trim(email)=$3 OR $3 is null);`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, statement, user.ID, user.Username, user.Email).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username,
		&u.Password, &u.Email, &u.JoinedAt, &u.IsVerified)

	return &u, err
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
							id, username, firstname, lastname, profil_picture
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

		err = rows.Scan(&u.ID, &u.Username, &u.Firstname, &u.Lastname, &u.ProfilPicture)
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

func (m *UserModel) Activate(u *User) error {
	const statement = `UPDATE users
						SET isverified=true
						WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, statement, u.ID)
	return err
}
