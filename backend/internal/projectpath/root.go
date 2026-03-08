package projectpath

import (
	"os"
	"path/filepath"
	"runtime"

	"backend/internal/apperr"
)

func ModuleRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", apperr.New(
			apperr.Internal,
			"projectpath.module_root_not_found",
			"failed to determine module root",
		)
	}

	current := filepath.Dir(filename)
	for {
		goModPath := filepath.Join(current, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", apperr.New(
				apperr.Internal,
				"projectpath.module_root_not_found",
				"failed to determine module root",
			)
		}

		current = parent
	}
}
