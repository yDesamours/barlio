package model

import (
	"barlio/internal/data"
	"database/sql"
)

type User struct {
	ID        int             `json:"id,omitempty"`
	Firstname data.DataString `json:"firstname,omitempty"`
	Lastname  data.DataString `json:"lastname,omitempty"`
	Username  data.DataString `json:"username,omitempty"`
	Password  data.DataString `json:"-"`
	Valid     bool            `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(u User) (int, error) {
	return u.ID, nil
}

func (m *UserModel) Get(id int) (User, error) {
	return User{}, nil
}

func (m *UserModel) Delete(id int) error {
	return nil
}

func (m *UserModel) Validate(id int) error {
	return nil
}
