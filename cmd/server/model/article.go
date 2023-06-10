package model

import (
	"barlio/internal/data"
	"barlio/internal/helper"
	"barlio/internal/metadata"
	"barlio/internal/validator"
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Article struct {
	ID         int
	Subject    data.String
	Categories ListArticleCategorie
	Title      data.String
	Content    data.String
	Gender     data.String
	CreatedAt  time.Time
	Author     User
	Tags       []data.String
}

type ArticleCategorie struct {
	Id   int `json:"id"`
	Name int `json:"name"`
}

type ListArticleCategorie []ArticleCategorie

type ArticleSearch struct {
	Subject    data.String
	Tags       []data.String
	Title      data.String
	CreatedAt  time.Time
	Author     data.String
	PageSize   int
	PageNumber int
}

type ArticleModel struct {
	DB *sql.DB
}

func (m *ArticleModel) Insert(art *Article) (int, error) {
	const statement = `INSERT INTO articles
							(subject, tags, title, content, id_author)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, statement, art.Subject, pq.Array(art.Tags),
		art.Title, art.Content, art.Author.ID).Scan(&art.ID)

	if err != nil {
		return 0, err
	}

	return art.ID, nil
}

func (m *ArticleModel) Update(art *Article) error {
	const statement = `UPDATE articles
						SET subject=$1, tags=$2, title=$3, content=$4
						WHERE id=$5
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, statement, art.Subject, pq.Array(art.Tags),
		art.Title, art.Content, art.ID).Scan(&art.ID)

	if err != nil {
		return err
	}
	return nil
}

func (m *ArticleModel) GetById(idArticle int) (*Article, error) {
	const statement = `SELECT 
							a.id, a.subject, a.content, a.tags, a.title, a.created_at
							u.id, u.username
						FROM articles a
						INNER JOIN users u ON u.id=a.id_author
						WHERE id=$1`

	var art Article
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, statement, idArticle).Scan(&art.ID, &art.Subject,
		&art.Content, pq.Array(&art.Tags), &art.Title, &art.CreatedAt, &art.Author.ID,
		&art.Author.Username)

	return &art, err
}

func (m *ArticleModel) Get(search *ArticleSearch) ([]Article, *metadata.Metadata, error) {
	var articles []Article
	var data metadata.Metadata
	const statement = `SELECT
							a.id, a.subject, a.tags, a.title, a.created_at
							u.id, u.username, count(a.id) OVER()
						FROM articles a
						INNER JOIN users u ON u.id=a.id_author
						WHERE 
							to_tsvector(a.subject) @@ plainto_tsquery($1) OR $1 =''
							AND u.username ilike $2 OR $2=''
							AND a.tags && $3 OR $3 =null
							LIMIT $4
							OFFSET $5`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, statement, search.Subject, search.Author,
		pq.Array(search.Tags), search.PageSize, search.PageNumber)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var art Article

		err = rows.Scan(&art.ID, &art.Subject, pq.Array(&art.Tags), &art.Title,
			&art.CreatedAt, &art.Author.ID, &art.Author.Username, &data.TotalResult)
		if err != nil {
			return nil, nil, err
		}
		articles = append(articles, art)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return articles, &data, nil
}

func (m *ArticleModel) Delete(idArticle, idAuthor int) error {
	const statement = `DELETE FROM articles
						WHERE 
							id=$1 AND id_author=$2
						RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return m.DB.QueryRowContext(ctx, statement, idArticle, idAuthor).Scan(&idArticle)
}

func (m *ArticleModel) Validate(art *Article, v *validator.Validator) {
	v.Check(len(art.Tags) > 2, "tags", "article must have at least 2 tags")
	v.Check(helper.IsASet(art.Tags), "tag", "article can't have a tag more than once")
	v.Check(art.Subject != "", "subject", "article subject is required")
	v.Check(art.Title != "", "title", "article title is required")
	v.Check(len(art.Content) > 200, "content", "article must be 200 character length")

}
