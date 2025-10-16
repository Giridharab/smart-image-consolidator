package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"os/exec"
)

func ScanDockerfiles(root string) []string {
	var dockerfiles []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "Dockerfile") {
			dockerfiles = append(dockerfiles, path)
		}
		return nil
	})
	return dockerfiles
}

func ReadDockerfile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func GetImageTag(dockerfilePath string) string {
	parts := strings.Split(dockerfilePath, "/")
	dir := parts[len(parts)-2] // parent directory of Dockerfile
	name := parts[len(parts)-1]

	tag := "pr-image-" + dir
	if strings.Contains(name, ".") {
		ext := strings.Split(name, ".")[1]
		tag += "-" + ext
	}
	// Replace invalid characters
	tag = strings.ToLower(tag)
	tag = strings.ReplaceAll(tag, "_", "-")
	tag = strings.ReplaceAll(tag, ".", "-")
	return tag
}

func BuildDockerImage(dockerfilePath string, tag string) error {
	cmd := exec.Command("docker", "build", "-f", dockerfilePath, "-t", tag, ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker build failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

