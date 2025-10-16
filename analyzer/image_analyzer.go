package analyzer

import (
	"strings"
)

func SuggestCanonicalBase(dockerfileContent string) string {
	lines := strings.Split(dockerfileContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(line), "FROM") {
			base := strings.TrimSpace(strings.Split(line, " ")[1])
			// Example: suggest slim variant if not already
			if !strings.Contains(base, "slim") {
				return base + "-slim"
			}
			return base
		}
	}
	return ""
}

