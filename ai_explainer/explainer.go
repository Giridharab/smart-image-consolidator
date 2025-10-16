package ai_explainer

import (
	"fmt"
	"smart-image-consolidator/scanner"
	"smart-image-consolidator/metrics"
	"smart-image-consolidator/analyzer"
)

func GenerateAIReport(image string, scanOutput string, dockerfileContent string) string {
	m := metrics.MeasurePerformanceAndCost(image)
	suggestions := analyzer.SuggestCanonicalBase(dockerfileContent)

	report := fmt.Sprintf("### 🚀 Smart Image Consolidator Report for `%s`\n\n", image)
	report += fmt.Sprintf("**Security Scan:**\n```\n%s\n```\n\n", scanOutput)
	report += fmt.Sprintf("**Performance:** CPU: %s, Memory: %s, Storage: %s\n", m.CPUUsage, m.MemoryUsage, m.Storage)
	report += fmt.Sprintf("**Estimated Daily Cost:** $%.2f\n", m.EstimatedCost)

	if !scanner.CheckChainGuardInstalled() {
		report += "\n**Recommendation:** ChainGuard is not enabled. Use ChainGuard to scan images for vulnerabilities.\n"
	}

	if len(suggestions) > 0 {
		report += "\n**Canonical Base Image Suggestions (Internal Registries):**\n"
		for _, s := range suggestions {
			report += fmt.Sprintf("- %s\n", s)
		}
	} else {
		report += "\n**Canonical Base Image Suggestion:** Use slim/minimal base images to reduce size and improve security.\n"
	}

	return report
}
