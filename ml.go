package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func CalculateMatch(url string, request Model) (int, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(http.MethodPost, url, &buf)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var s Similarity
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return 0, err
	}
	return s.Score, nil
}
