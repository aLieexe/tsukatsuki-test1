package mocks

import (
	"go-webserver/internal/models"
	"time"
)

var mockSnippet = models.Snippet{
	Id:        "snippet-123",
	Title:     "RIO RIO RIO RIO RIO RIO RIO RIO RIO",
	Content:   "RIO RIO RIO RIO RIO RIO ",
	CreatedAt: time.Now(),
	Expires:   time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (string, error) {
	return "snippet-1234", nil
}

func (m *SnippetModel) Get(id string) (models.Snippet, error) {
	switch id {
	case "snippet-123":
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}
func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
