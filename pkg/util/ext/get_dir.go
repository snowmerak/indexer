package ext

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetDirectoryName(path string) (string, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	split := strings.Split(filepath.Dir(dir), string(os.PathSeparator))
	if len(split) > 0 {
		dir = split[len(split)-1]
	}
	if dir == "" {
		dir = "root_directory"
	}

	return dir, nil
}
