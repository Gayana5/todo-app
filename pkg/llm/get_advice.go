package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *OllamaClient) GetAdvice(title, description string) (string, error) {
	reqBody := map[string]interface{}{
		"model": "gemma:2b", // Используем явное указание версии
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("Give me advice on how to achieve the goal: '%s'. Description: '%s'. Write 5 tips in json format.", title, description),
			},
		},
		"stream": false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post("http://ollama:11434/api/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем код статуса
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("Ollama error: %s", result.Error)
	}

	return result.Message.Content, nil
}
