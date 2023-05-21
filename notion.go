package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Notion struct {
	http       *http.Client
	DatabaseID string
	Token      string
}

func NewNotion(databaseID string, token string) *Notion {
	return &Notion{
		DatabaseID: databaseID,
		Token:      token,
		http:       &http.Client{},
	}
}

func (n *Notion) Send(amount float64, category string) error {
	payload := NotionPostPayload{
		Parent: Parent{
			Type:       "database_id",
			DatabaseID: n.DatabaseID,
		},
		Properties: map[string]interface{}{
			"title":    NewTitleProperty("test"),
			"amount":   NewAmountProperty(amount),
			"category": NewCategoryProperty(category),
		},
	}
	payloadBytes, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", "https://api.notion.com/v1/pages", payloadReader)
	req.Header = http.Header{
		"Authorization":  []string{"Bearer " + n.Token},
		"Content-Type":   []string{"application/json"},
		"Notion-Version": []string{"2022-06-28"},
	}
	resp, err := n.http.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	respBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&respBody)
	fmt.Printf("%v \n", respBody)
	return nil
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
