package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type LLM struct {
	model  string
	llmURL string
}

var response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Think  bool   `json:"think"`
}

func New(model, llmURL, port string) *LLM {
	return &LLM{
		model:  model,
		llmURL: llmURL,
	}
}

func (l *LLM) GetRating(q string) (int, error) {
	if q == "" {
		return 0, fmt.Errorf("empty query")
	}

	prompt := fmt.Sprintf("Оцени отзыв. Верни результат в виде одного целого числа. 1 - негативный, 2 - позитивный: %s", q)
	payload := LLMRequest{
		Model:  l.model,
		Prompt: prompt,
		Stream: false,
		Think:  true,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования в JSON: %v", err)
	}

	req, err := http.NewRequest("POST", l.llmURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("запрос завершился с ошибкой %d: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	if !response.Done {
		return 0, fmt.Errorf("ответ не завершен")
	}

	rating, err := strconv.Atoi(response.Response)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования в число: %v", err)
	}
	return rating, nil
}
