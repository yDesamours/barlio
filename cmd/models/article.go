package models

import (
	"database/sql"
	"time"
)

type Article struct {
	ID        int       `json:"id"`
	Subject   string    `json:"subject,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Author    User      `json:"author"`
}

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) Insert(art Article) (int, error) {
	const statement = `INSERT INTO articles
							(subject, tags, title, content, id_author)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING id`

	return art.ID, nil
}

func (m *ArticleModel) Update(art Article) error {
	return nil
}

func (m *ArticleModel) Get(id int) (*Article, error) {
	return &Article{}, nil
}

func (art *ArticleModel) Latest() (*Article, error) {
	return &Article{}, nil
}

func (art *ArticleModel) Delete(id int) error {
	return nil
}
