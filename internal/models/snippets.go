package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	Conn *pgx.Conn
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	var id int

	stmt := `insert into snippets (title, content, created, expires)
	values($1, $2, now(), now() + interval '1 day' * $3 ) RETURNING id`
	err := m.Conn.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return 0, err
		}
	}

	return id, nil

}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := &Snippet{}

	err := m.Conn.QueryRow(context.Background(), `select id, title, content, created, expires
		from snippets where expires > now() and id = $1`, id).Scan(&s.ID, &s.Title,
		&s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
