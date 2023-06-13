package model

import (
	"database/sql"
	"time"
)

const (
	VerificationScope = "verification"
)

type Token struct {
	Userid    int
	Token     string
	Scope     string
	Hash      string
	ExpiretAt time.Time
}

type TokenModel struct {
	DB *sql.DB
}

func (t *TokenModel) Insert(token *Token) error {
	return nil
}

func (t *TokenModel) DeleteExpired() error {
	return nil
}

func (t *TokenModel) DeleteForUser(userId int, scope string) error {
	return nil
}
