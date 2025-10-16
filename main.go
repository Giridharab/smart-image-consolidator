package main

import (
	"fmt"
	"smart-image-consolidator/scanner"
	"smart-image-consolidator/ai_explainer"
	"smart-image-consolidator/ci"
)

func main() {
	dockerfiles := scanner.ScanDockerfiles(".")
	fmt.Printf("Found %d Dockerfiles in PR\n", len(dockerfiles))

	for _, df := range dockerfiles {
		tag := scanner.GetImageTag(df)
		content := scanner.ReadDockerfile(df)

		err := scanner.BuildDockerImage(df, tag)
		var scanResult string
		if err != nil {
			scanResult = fmt.Sprintf("Image build failed: %v", err)
		} else if scanner.CheckChainGuardInstalled() {
			result, err := scanner.ScanImageWithChainGuard(tag)
			if err != nil {
				scanResult = fmt.Sprintf("ChainGuard scan failed: %v", err)
			} else {
				scanResult = result
			}
		} else {
			scanResult = "ChainGuard not installed. Security scan skipped."
		}

		report := ai_explainer.GenerateAIReport(tag, scanResult, content)
		ci.CommentOnPR(report)
	}

	fmt.Println("✅ Smart Image Consolidator completed for all Dockerfiles")
}
