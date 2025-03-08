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
	// SQL query to grab the most recently created snippets
	stmt := `select id, title, content, created, expires from snippets
		where expires > now() order by id desc limit 10`

	rows, err := m.Conn.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	// makes sure resultset is close before we return the results. close after error check
	// this way we are sure the database connection is no longer open
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// when rows.Next() finishes, call rows.Err() to retrieve any error that
	// has was encountered during the iteration. Never assume a successful
	// iteration was completed over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
