package config

import (
	"errors"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"backend/internal/apperr"
	"backend/internal/projectpath"

	"github.com/joho/godotenv"
)

const defaultPort = 8698
const defaultDevelopmentOrigin = "http://localhost:3000"

const (
	appEnvVar            = "TONARIS_ENV"
	portEnvVar           = "PORT"
	allowedOriginsEnvVar = "CORS_ALLOWED_ORIGINS"
)

type RunEnvironment string

const (
	Development RunEnvironment = "development"
	Production  RunEnvironment = "production"
)

type Config struct {
	Environment    RunEnvironment
	Port           int
	AllowedOrigins []string
}

func Load() (Config, error) {
	moduleRoot, err := projectpath.ModuleRoot()
	if err != nil {
		return Config{}, apperr.Wrap(
			err,
			apperr.Internal,
			"projectpath.module_root_not_found",
			"failed to load configuration",
		)
	}

	return loadFromEnv(envMap(os.Environ()), filepath.Join(moduleRoot, ".env"))
}

func loadFromEnv(processEnv map[string]string, dotenvPath string) (Config, error) {
	mergedEnv := cloneEnv(processEnv)

	// Local dotenv values provide development defaults, but the process
	// environment keeps highest precedence in every environment.
	if shouldLoadDotenv(mergedEnv[appEnvVar]) {
		dotenvEnv, err := readDotenv(dotenvPath)
		if err != nil {
			return Config{}, err
		}

		for key, value := range dotenvEnv {
			if _, exists := mergedEnv[key]; !exists {
				mergedEnv[key] = value
			}
		}
	}

	return parseEnv(mergedEnv)
}

func parseEnv(values map[string]string) (Config, error) {
	environment, err := parseEnvironment(values[appEnvVar])
	if err != nil {
		return Config{}, err
	}

	port, err := parsePort(values[portEnvVar])
	if err != nil {
		return Config{}, err
	}

	allowedOrigins, err := parseAllowedOrigins(environment, values[allowedOriginsEnvVar])
	if err != nil {
		return Config{}, err
	}

	return Config{
		Environment:    environment,
		Port:           port,
		AllowedOrigins: allowedOrigins,
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

func parseAllowedOrigins(environment RunEnvironment, raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		if environment == Production {
			return nil, apperr.New(
				apperr.InvalidArgument,
				"config.missing_cors_allowed_origins",
				"missing CORS allowed origins configuration",
			)
		}

		return []string{defaultDevelopmentOrigin}, nil
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			return nil, apperr.New(
				apperr.InvalidArgument,
				"config.invalid_cors_allowed_origins",
				"invalid CORS allowed origins configuration",
			)
		}

		origins = append(origins, origin)
	}

	return origins, nil
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
	maps.Copy(cloned, values)

	return cloned
}
