package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	"net/http"
)

var (
	ollamaURL string
	model     string
	port      string
)

var response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func main() {
	flag.StringVar(&port, "p", "8080", "port")
	flag.StringVar(&ollamaURL, "u", "http://localhost:11434/api/generate", "ollama url")
	flag.StringVar(&model, "m", "deepseek-r1:1.5b", "model")
	flag.Parse()
	http.HandleFunc("/", query)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func query(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	q := params.Get("q")
	res, err := ollamaRequest(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(res))
}

func ollamaRequest(q string) (string, error) {
	payload := map[string]interface{}{
		"model":  model,
		"prompt": q,
		"stream": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if !response.Done {
		return "", fmt.Errorf("incomplete response from Ollama")
	}

	return response.Response, nil
}
