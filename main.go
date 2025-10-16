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

		// Optional: still build image if you want metrics
		err := scanner.BuildDockerImage(df, tag)
		if err != nil {
			fmt.Printf("Image build failed: %v\n", err)
		}

		// Skip ChainGuard scan
		scanResult := "Security scan skipped (ChainGuard not installed)."

		// Generate report with canonical base suggestions
		report := ai_explainer.GenerateAIReport(tag, scanResult, content)

		// Post PR comment
		ci.CommentOnPR(report)
	}

	fmt.Println("✅ Smart Image Consolidator completed for all Dockerfiles")
}
