package model

import (
	"barlio/internal/validator"
	"database/sql"
	"testing"
)

func TestInsertExistingUser(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://barlio:barliopass@localhost:5432/barlio?sslmode=disable")
	if err != nil {
		t.Fail()
	}

	err = db.Ping()
	if err != nil {
		t.Fail()
	}

	userModel := UserModel{DB: db}
	user := User{Username: "ydesamours", Email: "yveltamours@gmail.com", Password: "pass"}
	v := validator.New()
	userModel.ValidateUser(&user, v)
	if v.Valid() {
		t.Fail()
	}

}
