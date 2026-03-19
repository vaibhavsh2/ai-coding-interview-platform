package service

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const baseURL = "https://judge0-ce.p.rapidapi.com"

type JudgeRequest struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin"`
}

type JudgeResponse struct {
	Token string `json:"token"`
}

type JudgeResult struct {
	Stdout string `json:"stdout"`
	Status struct {
		Description string `json:"description"`
	} `json:"status"`
}

// Step 1: Send code to Judge0
func SubmitToJudge(code string, input string) (string, error) {

	body := JudgeRequest{
		LanguageID: 62, // Java
		SourceCode: code,
		Stdin:      input,
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		baseURL+"/submissions?base64_encoded=false&wait=false",
		bytes.NewBuffer(jsonData),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", "8dc20015a8msh206bdb51e75d11bp18232bjsn6c421d0d290e")
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result JudgeResponse
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Token, nil
}

// Step 2: Get result from Judge0
func GetJudgeResult(token string) (JudgeResult, error) {

	req, _ := http.NewRequest(
		"GET",
		baseURL+"/submissions/"+token+"?base64_encoded=false",
		nil,
	)

	req.Header.Set("X-RapidAPI-Key", "8dc20015a8msh206bdb51e75d11bp18232bjsn6c421d0d290e")
	req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return JudgeResult{}, err
	}

	defer resp.Body.Close()

	var result JudgeResult
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
