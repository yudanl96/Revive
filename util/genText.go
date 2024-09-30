package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GenerateText(prompt string, length int) (string, error) {
	url := "https://api-inference.huggingface.co/models/mistralai/Mixtral-8x7B-Instruct-v0.1" // Change the model as needed
	payload := map[string]interface{}{
		"inputs":     prompt,
		"parameters": map[string]int{"max_new_tokens": length},
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer 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") // Add your API key
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Change the variable type to slice of interfaces to handle array response
	var result []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Convert the response to string or handle it as needed
	if len(result) > 0 {
		if output, ok := result[0].(map[string]interface{}); ok {
			if text, ok := output["generated_text"].(string); ok {
				return text[len(prompt):], nil
			}
		}
	}

	return "", fmt.Errorf("unexpected response format")
}
