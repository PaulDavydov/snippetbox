package models

import (
	"time"

	"github.com/jackc/pgx/v5"
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
	return 0, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
