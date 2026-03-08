package projectpath

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func ModuleRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("determine caller path")
	}

	current := filepath.Dir(filename)
	for {
		goModPath := filepath.Join(current, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("find module root from %s", filename)
		}

		current = parent
	}
}
