package analyzer

import (
	"strings"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"fmt"
)

type CanonicalBases map[string][]string

func SuggestCanonicalBase(dockerfileContent string) []string {
	data, err := ioutil.ReadFile("configs/canonical_bases.yaml")
	if err != nil {
		fmt.Println("Failed to read canonical_bases.yaml:", err)
		return nil
	}
	var bases CanonicalBases
	err = yaml.Unmarshal(data, &bases)
	if err != nil {
		fmt.Println("Failed to parse canonical_bases.yaml:", err)
		return nil
	}

	lines := strings.Split(dockerfileContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToUpper(line), "FROM") {
			currentBase := strings.TrimSpace(strings.Split(line, " ")[1])
			for key, regs := range bases {
				if strings.Contains(currentBase, key) {
					return regs
				}
			}
			// fallback suggestion: all internal registries
			var suggestions []string
			for _, regs := range bases {
				suggestions = append(suggestions, regs...)
			}
			return suggestions
		}
	}
	return nil
}
