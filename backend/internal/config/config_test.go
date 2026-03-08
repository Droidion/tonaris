package config

import (
	"os"
	"path/filepath"
	"testing"

	"backend/internal/apperr"
)

func TestLoadFromEnvUsesDefaults(t *testing.T) {
	t.Parallel()

	cfg, err := loadFromEnv(map[string]string{}, filepath.Join(t.TempDir(), ".env"))
	if err != nil {
		t.Fatalf("loadFromEnv returned error: %v", err)
	}

	if cfg.Environment != Development {
		t.Fatalf("expected development environment, got %q", cfg.Environment)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("expected port %d, got %d", defaultPort, cfg.Port)
	}
}

func TestLoadFromEnvRejectsInvalidPort(t *testing.T) {
	t.Parallel()

	_, err := loadFromEnv(map[string]string{portEnvVar: "invalid"}, filepath.Join(t.TempDir(), ".env"))
	if err == nil {
		t.Fatal("expected invalid port error")
	}

	appErr, ok := apperr.As(err)
	if !ok {
		t.Fatalf("expected app error, got %T", err)
	}

	if appErr.Kind != apperr.InvalidArgument {
		t.Fatalf("expected invalid argument kind, got %q", appErr.Kind)
	}

	if appErr.Code != "config.invalid_port" {
		t.Fatalf("expected config.invalid_port code, got %q", appErr.Code)
	}
}

func TestLoadFromEnvRejectsInvalidEnvironment(t *testing.T) {
	t.Parallel()

	_, err := loadFromEnv(map[string]string{appEnvVar: "staging"}, filepath.Join(t.TempDir(), ".env"))
	if err == nil {
		t.Fatal("expected invalid environment error")
	}

	appErr, ok := apperr.As(err)
	if !ok {
		t.Fatalf("expected app error, got %T", err)
	}

	if appErr.Code != "config.invalid_environment" {
		t.Fatalf("expected config.invalid_environment code, got %q", appErr.Code)
	}
}

func TestLoadFromEnvReadsDotenvInDevelopment(t *testing.T) {
	t.Parallel()

	dotenvPath := writeDotenv(t, "TONARIS_ENV=development\nPORT=9001\n")

	cfg, err := loadFromEnv(map[string]string{}, dotenvPath)
	if err != nil {
		t.Fatalf("loadFromEnv returned error: %v", err)
	}

	if cfg.Environment != Development {
		t.Fatalf("expected development environment, got %q", cfg.Environment)
	}

	if cfg.Port != 9001 {
		t.Fatalf("expected port 9001, got %d", cfg.Port)
	}
}

func TestLoadFromEnvPrefersProcessEnvOverDotenv(t *testing.T) {
	t.Parallel()

	dotenvPath := writeDotenv(t, "TONARIS_ENV=development\nPORT=9001\n")

	cfg, err := loadFromEnv(map[string]string{
		appEnvVar:  "development",
		portEnvVar: "9100",
	}, dotenvPath)
	if err != nil {
		t.Fatalf("loadFromEnv returned error: %v", err)
	}

	if cfg.Environment != Development {
		t.Fatalf("expected development environment, got %q", cfg.Environment)
	}

	if cfg.Port != 9100 {
		t.Fatalf("expected port 9100, got %d", cfg.Port)
	}
}

func TestLoadFromEnvSkipsDotenvInProduction(t *testing.T) {
	t.Parallel()

	dotenvPath := writeDotenv(t, "TONARIS_ENV=development\nPORT=9001\n")

	cfg, err := loadFromEnv(map[string]string{appEnvVar: "production"}, dotenvPath)
	if err != nil {
		t.Fatalf("loadFromEnv returned error: %v", err)
	}

	if cfg.Environment != Production {
		t.Fatalf("expected production environment, got %q", cfg.Environment)
	}

	if cfg.Port != defaultPort {
		t.Fatalf("expected port %d, got %d", defaultPort, cfg.Port)
	}
}

func writeDotenv(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write dotenv: %v", err)
	}

	return path
}
