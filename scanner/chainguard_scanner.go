package scanner

import (
	"fmt"
	"os/exec"
)

func CheckChainGuardInstalled() bool {
	_, err := exec.LookPath("chainguard")
	return err == nil
}

func ScanImageWithChainGuard(image string) (string, error) {
	cmd := exec.Command("chainguard", "scan", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ChainGuard scan failed: %v\nOutput: %s", err, string(output))
	}
	return string(output), nil
}

