package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	//placeholder
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Snippet struct {
	Id        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Expires   time.Time `json:"expires" db:"expires"`
}

type SnippetRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Expires int    `json:"expires"`
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (string, error)
	Get(id string) (Snippet, error)
	Latest() ([]Snippet, error)
}

type SnippetModel struct {
	Pool *pgxpool.Pool
}

func (m *SnippetModel) ParseRequest(reqBody string) (SnippetRequest, error) {
	var parsedRequest SnippetRequest
	err := json.Unmarshal([]byte(reqBody), &parsedRequest)
	if err != nil {
		return SnippetRequest{}, fmt.Errorf("models: balls")
	}

	return parsedRequest, nil
}
func (m *SnippetModel) Insert(title string, content string, expires int) (string, error) {
	id, err := gonanoid.New(16)
	id = fmt.Sprint("snippet-", id)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO snippets(id, title, content, created_at, expires) VALUES
	(@id, @title, @content, @createdAt, @expires)`

	args := pgx.NamedArgs{
		"id":        id,
		"title":     title,
		"content":   content,
		"createdAt": time.Now(),
		"expires":   time.Now().AddDate(0, 0, expires),
	}

	commandTag, err := m.Pool.Exec(context.Background(), query, args)

	if err != nil {
		return "", err
	}

	if commandTag.RowsAffected() != 1 {
		return "", err
	}

	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id string) (Snippet, error) {
	query := `SELECT * FROM snippets WHERE expires > CURRENT_TIMESTAMP AND id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := m.Pool.Query(context.Background(), query, args)
	if err != nil {
		return Snippet{}, err
	}

	snippet, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Snippet])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return snippet, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	query := `SELECT * FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY created_at DESC LIMIT 5`
	rows, err := m.Pool.Query(context.Background(), query)
	if err != nil {
		return []Snippet{}, err
	}

	snippets, err := pgx.CollectRows(rows, pgx.RowToStructByName[Snippet])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []Snippet{}, ErrNoRecord
		}
		return []Snippet{}, err
	}

	return snippets, nil
}

func (m *SnippetModel) Check() {
	query := `SELECT data FROM sessions`
	rows, err := m.Pool.Query(context.Background(), query)

	if err != nil {
		fmt.Println(err)
	}

	res, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) ([]byte, error) {
		var data []byte
		err := row.Scan(&data)
		return data, err
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(res)
	fmt.Println()
	fmt.Println()

}
