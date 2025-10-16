package main

import (
	"fmt"
	"smart-image-consolidator/scanner"
	"smart-image-consolidator/metrics"
	"smart-image-consolidator/ai_explainer"
	"smart-image-consolidator/ci"
)

func main() {
	dockerfiles := scanner.ScanDockerfiles(".")
	fmt.Printf("Found %d Dockerfiles in PR\n", len(dockerfiles))

	for _, df := range dockerfiles {
		tag := scanner.GetImageTag(df)
		content := scanner.ReadDockerfile(df)

		// Measure real-time metrics for the image
		perfMetrics, err := metrics.MeasurePerformanceAndCost(tag)
		var report string
		if err != nil {
			report = fmt.Sprintf("⚠️ Failed to measure performance for image %s: %v", tag, err)
		} else {
			report = fmt.Sprintf("Docker Image: %s\nCPU: %s\nMemory: %s\nStorage: %s\nEstimated Cost: $%.2f",
				tag, perfMetrics.CPUUsage, perfMetrics.MemoryUsage, perfMetrics.Storage, perfMetrics.EstimatedCost)
		}

		// Optionally, add AI suggestions for canonical image
		aiReport := ai_explainer.GenerateAIReport(tag, report, content)
		finalReport := fmt.Sprintf("%s\n\n%s", report, aiReport)

		// Post comment to PR
		ci.CommentOnPR(finalReport)
	}

	fmt.Println("✅ Smart Image Consolidator completed for all Dockerfiles")
}
