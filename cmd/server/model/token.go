package model

import (
	"context"
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
	const statement = `INSERT INTO tokens
							(userid, scope, hash, expire_at)
						VALUES($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, statement, token.Userid, token.Scope, token.Hash, token.ExpiretAt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenModel) DeleteExpired() error {
	const statement = `DELETE FROM tokens 
						WHERE expired_at < $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, statement, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenModel) DeleteForUser(userId int, scope string) error {
	const statement = `DELETE FROM tokens 
						WHERE userid=$1 AND scope=$2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, statement, userId, scope)
	if err != nil {
		return err
	}
	return nil
}
