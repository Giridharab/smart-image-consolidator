package ai_explainer

import (
	"fmt"
	"smart-image-consolidator/scanner"
	"smart-image-consolidator/metrics"
)

func GenerateAIReport(image string, scanOutput string) string {
	m := metrics.MeasurePerformanceAndCost(image)

	report := fmt.Sprintf(
		"### 🚀 Smart Image Consolidator Report for `%s`\n\n", image)
	report += fmt.Sprintf("**Security Scan:**\n```\n%s\n```\n\n", scanOutput)
	report += fmt.Sprintf("**Performance:** CPU: %s, Memory: %s, Storage: %s\n", m.CPUUsage, m.MemoryUsage, m.Storage)
	report += fmt.Sprintf("**Estimated Daily Cost:** $%.2f\n", m.EstimatedCost)

	if !scanner.CheckChainGuardInstalled() {
		report += "\n**Recommendation:** ChainGuard is not enabled. Use ChainGuard to scan images for vulnerabilities.\n"
	}

	report += "\n**Canonical Base Image Suggestion:** Use slim/minimal base images to reduce size and improve security.\n"
	return report
}

