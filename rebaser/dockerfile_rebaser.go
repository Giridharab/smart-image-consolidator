package rebaser

import (
	"io/ioutil"
	"strings"
)

func RebaseDockerfile(path string, newBase string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.ToUpper(line), "FROM") {
			lines[i] = "FROM " + newBase
			break
		}
	}
	newContent := strings.Join(lines, "\n")
	return ioutil.WriteFile(path, []byte(newContent), 0644)
}

