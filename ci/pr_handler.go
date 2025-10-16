package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SetPRStatus(report string, conclusion string) {
	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY")
	commitSHA := os.Getenv("COMMIT_SHA")

	if token == "" || commitSHA == "" || repo == "" {
		fmt.Println("Missing environment variables for PR status")
		return
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/check-runs", repo)
	body := map[string]interface{}{
		"name":       "Smart Image Consolidator",
		"head_sha":   commitSHA,
		"status":     "completed",
		"conclusion": conclusion, // "success" | "failure" | "neutral"
		"output": map[string]string{
			"title":   "Smart Image Consolidator Analysis",
			"summary": report,
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to post PR status:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("PR status posted with response code:", resp.StatusCode)
}
