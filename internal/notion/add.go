package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (n *Notion) Add(amount float64, category, title string) error {
	props := map[string]interface{}{
		"amount":   NewAmountProperty(amount),
		"category": NewCategoryProperty(category),
	}

	if title != "" {
		props["title"] = NewTitleProperty(title)
	} else {
		props["title"] = NewTitleProperty(fmt.Sprintf("฿ %.2f spend with %s", amount, category))
	}
	payload := NotionPostPayload{
		Parent: Parent{
			Type:       "database_id",
			DatabaseID: n.DatabaseID,
		},
		Properties: props,
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
		n.logger.Error(
			"Request failed",
			slog.String("errorMessage", err.Error()),
			slog.String("requestUrl", req.URL.String()),
			slog.String("requestMethod", req.Method),
		)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		responseError := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseError)
		n.logger.Error(
			"Request failed",
			slog.String("requestUrl", req.URL.String()),
			slog.String("requestMethod", req.Method),
			slog.Int("statusCode", resp.StatusCode),
			slog.Any("responseBody", responseError),
		)
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	respBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&respBody)
	return nil
}
