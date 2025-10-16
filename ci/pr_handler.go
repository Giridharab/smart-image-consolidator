package ci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func CommentOnPR(report string) {
	token := os.Getenv("GITHUB_TOKEN")
	prNumber := os.Getenv("GITHUB_PR_NUMBER")
	repo := os.Getenv("GITHUB_REPOSITORY")
	commitSHA := os.Getenv("COMMIT_SHA")

	if token == "" || prNumber == "" || repo == "" || commitSHA == "" {
		fmt.Println("Missing environment variables")
		return
	}

	// Create PR comment
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, prNumber)
	body, _ := json.Marshal(map[string]string{"body": report})
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)

	// Optional: set PR status
	statusURL := fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", repo, commitSHA)
	statusBody := map[string]string{
		"state":       "success",
		"description": "Smart Image Consolidator: Analysis complete",
		"context":     "Smart Image Consolidator",
	}
	jsonStatus, _ := json.Marshal(statusBody)
	req2, _ := http.NewRequest("POST", statusURL, bytes.NewBuffer(jsonStatus))
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/json")
	client.Do(req2)
}

