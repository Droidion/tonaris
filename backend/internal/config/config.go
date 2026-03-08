package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"backend/internal/apperr"
	"github.com/joho/godotenv"

	"backend/internal/projectpath"
)

const defaultPort = 8698

const (
	appEnvVar  = "TONARIS_ENV"
	portEnvVar = "PORT"
)

type RunEnvironment string

const (
	Development RunEnvironment = "development"
	Production  RunEnvironment = "production"
)

type AppConfig struct {
	Environment RunEnvironment
	Port        int
}

func Load() (AppConfig, error) {
	moduleRoot, err := projectpath.ModuleRoot()
	if err != nil {
		return AppConfig{}, apperr.Wrap(
			err,
			apperr.Internal,
			"projectpath.module_root_not_found",
			"failed to load configuration",
		)
	}

	return loadFromEnv(envMap(os.Environ()), filepath.Join(moduleRoot, ".env"))
}

func loadFromEnv(processEnv map[string]string, dotenvPath string) (AppConfig, error) {
	mergedEnv := cloneEnv(processEnv)

	if shouldLoadDotenv(mergedEnv[appEnvVar]) {
		dotenvEnv, err := readDotenv(dotenvPath)
		if err != nil {
			return AppConfig{}, err
		}

		for key, value := range dotenvEnv {
			if _, exists := mergedEnv[key]; !exists {
				mergedEnv[key] = value
			}
		}
	}

	return parseEnv(mergedEnv)
}

func parseEnv(values map[string]string) (AppConfig, error) {
	environment, err := parseEnvironment(values[appEnvVar])
	if err != nil {
		return AppConfig{}, err
	}

	port, err := parsePort(values[portEnvVar])
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		Environment: environment,
		Port:        port,
	}, nil
}

func parseEnvironment(raw string) (RunEnvironment, error) {
	if raw == "" {
		return Development, nil
	}

	environment := RunEnvironment(raw)
	switch environment {
	case Development, Production:
		return environment, nil
	default:
		return "", apperr.New(
			apperr.InvalidArgument,
			"config.invalid_environment",
			"invalid environment configuration",
		)
	}
}

func parsePort(raw string) (int, error) {
	if raw == "" {
		return defaultPort, nil
	}

	port, err := strconv.ParseUint(raw, 10, 16)
	if err != nil {
		return 0, apperr.Wrap(
			err,
			apperr.InvalidArgument,
			"config.invalid_port",
			"invalid port configuration",
		)
	}

	return int(port), nil
}

func shouldLoadDotenv(rawEnvironment string) bool {
	return rawEnvironment != string(Production)
}

func readDotenv(dotenvPath string) (map[string]string, error) {
	_, err := os.Stat(dotenvPath)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, apperr.Wrap(
			err,
			apperr.Internal,
			"config.dotenv_read_failed",
			"failed to read dotenv file",
		)
	}

	values, err := godotenv.Read(dotenvPath)
	if err != nil {
		return nil, apperr.Wrap(
			err,
			apperr.Internal,
			"config.dotenv_read_failed",
			"failed to read dotenv file",
		)
	}

	return values, nil
}

func envMap(values []string) map[string]string {
	env := make(map[string]string, len(values))
	for _, entry := range values {
		key, value, found := strings.Cut(entry, "=")
		if found {
			env[key] = value
		}
	}

	return env
}

func cloneEnv(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}
