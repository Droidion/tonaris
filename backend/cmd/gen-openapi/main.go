package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"backend/internal/api"
	"backend/internal/projectpath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "generate schema: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := api.New(api.Dependencies{})
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(app.API.OpenAPI(), "", "  ")
	if err != nil {
		return fmt.Errorf("marshal openapi: %w", err)
	}

	moduleRoot, err := projectpath.ModuleRoot()
	if err != nil {
		return err
	}

	outputPath := filepath.Join(moduleRoot, "..", "shared", "openapi.json")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create shared directory: %w", err)
	}

	jsonBytes = append(jsonBytes, '\n')
	if err := os.WriteFile(outputPath, jsonBytes, 0o644); err != nil {
		return fmt.Errorf("write schema: %w", err)
	}

	fmt.Printf("Schema written to %s\n", outputPath)

	return nil
}
