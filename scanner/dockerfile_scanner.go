package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"os/exec"
	"crypto/sha1"
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
	var dir string
	if len(parts) > 1 {
		dir = parts[len(parts)-2]
	} else {
		dir = "root"
	}

	name := parts[len(parts)-1]
	tagBase := "pr-image-" + dir

	// add suffix if Dockerfile has a name like Dockerfile.python
	if strings.Contains(name, ".") {
		ext := strings.Split(name, ".")[1]
		tagBase += "-" + ext
	}

	// add short hash to guarantee uniqueness
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(dockerfilePath)))[:6]

	tag := fmt.Sprintf("%s-%s", tagBase, hash)

	// clean up tag (Docker tags must be lowercase, no spaces)
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

