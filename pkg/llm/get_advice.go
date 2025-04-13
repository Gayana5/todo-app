package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *LLMClient) SendToHyperbolic(title, description string) (string, error) {
	url := "https://api.hyperbolic.xyz/v1/chat/completions"
	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJnc2RhbGxha2lhbkBnbWFpbC5jb20iLCJpYXQiOjE3NDQxMDY4NzZ9.O9hCJjACWOmadPsFxpTUbT8CV2ctzqov5XFNJ2b5K1s"

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("Give me advice on how to achieve the goal, 4 sentences max: '%s'. Description: '%s'. Answer in the language of the goal and description.", title, description),
			},
		},
		"model":       "Qwen/Qwen2.5-72B-Instruct",
		"max_tokens":  512,
		"temperature": 0.7,
		"top_p":       0.9,
		"stream":      false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}
