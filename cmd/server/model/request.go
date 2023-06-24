package model

import (
	"context"
	"database/sql"
	"time"
)

type Request struct {
	UserId    int
	Scope     string
	CreatedAt time.Time
	Data      []byte
	ExpireAt  time.Time
}

type RequestModel struct {
	db *sql.DB
}

func (m *RequestModel) New(db *sql.DB) {
	m.db = db
}

func (m *RequestModel) Insert(r *Request) error {
	const statement = `INSERT INTO requests (userid, scope, createdat, data, expireat)
						VALUES ($1, $2, $3, $4, $5)`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, statement, r.UserId, r.Scope, time.Now(), r.Data, r.ExpireAt)
	return err
}

func (m *RequestModel) GetForUser(userId int, scope string) (*Request, error) {
	var r Request
	const statement = `SELECT scope, createdat, data, expireat
						FROM requests
						WHERE userid=$1 AND scope=$2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := m.db.QueryRowContext(ctx, statement, userId, scope).Scan(&r.Scope, &r.CreatedAt,
		&r.Data, &r.ExpireAt)

	return &r, err
}

func (m *RequestModel) DeleteForUser(userId int, scope string) (*Request, error) {
	var r Request
	const statement = `DELETE 
						FROM requests
						WHERE userid=$1 AND type=$2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, statement, userId, scope)

	return &r, err
}
