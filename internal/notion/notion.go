package notion

import (
	"log/slog"
	"net/http"
)

type Notion struct {
	http       *http.Client
	DatabaseID string
	Token      string
	logger     *slog.Logger
}

func NewNotion(databaseID string, token string, logger *slog.Logger) *Notion {
	logger.Info("Notion client created")
	return &Notion{
		DatabaseID: databaseID,
		Token:      token,
		http:       &http.Client{},
		logger:     logger,
	}
}

type Parent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

type NotionPostPayload struct {
	Parent     Parent                 `json:"parent"`
	Properties map[string]interface{} `json:"properties"`
}

func NewTitleProperty(title string) map[string]interface{} {
	return map[string]interface{}{
		"type": "title",
		"title": []map[string]interface{}{
			{
				"type": "text",
				"text": map[string]interface{}{
					"content": title,
				},
			},
		},
	}
}

func NewAmountProperty(amount float64) map[string]interface{} {
	return map[string]interface{}{
		"type":   "number",
		"number": amount,
	}
}

func NewCategoryProperty(category string) map[string]interface{} {
	return map[string]interface{}{
		"type": "select",
		"select": map[string]interface{}{
			"name": category,
		},
	}
}
