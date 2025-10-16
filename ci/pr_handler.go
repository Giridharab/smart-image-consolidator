package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// CommentOnPR posts a comment to the pull request
func CommentOnPR(report string) {
	prNumber := os.Getenv("GITHUB_PR_NUMBER")
	repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("GITHUB_TOKEN")

	if prNumber == "" || repo == "" || token == "" {
		fmt.Println("Missing environment variables for PR comment")
		return
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, prNumber)

	body := map[string]string{
		"body": report,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to post comment:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("PR comment posted with response code:", resp.StatusCode)
}