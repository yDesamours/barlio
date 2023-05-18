package models

import "database/sql"

type User struct {
	ID        int    `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"-"`
	Valid     bool   `json:"-"`
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
