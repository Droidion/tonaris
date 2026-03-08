package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type Options struct {
	Environment string
	Level       string
	Format      string
	AddSource   bool
	ServiceName string
	Output      io.Writer
}

func New(opts Options) (*slog.Logger, error) {
	environment := strings.ToLower(strings.TrimSpace(opts.Environment))
	if environment == "" {
		environment = "development"
	}

	level, err := parseLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	format, err := resolveFormat(opts.Format, environment)
	if err != nil {
		return nil, err
	}

	output := opts.Output
	if output == nil {
		output = os.Stdout
	}

	serviceName := strings.TrimSpace(opts.ServiceName)
	if serviceName == "" {
		serviceName = "backend"
	}

	handler, err := newHandler(output, format, level, opts.AddSource)
	if err != nil {
		return nil, err
	}

	return slog.New(handler).With(
		"service", serviceName,
		"env", environment,
	), nil
}

func newHandler(output io.Writer, format string, level slog.Level, addSource bool) (slog.Handler, error) {
	switch format {
	case "json":
		return slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource:   addSource,
			Level:       level,
			ReplaceAttr: replaceAttr,
		}), nil
	case "pretty":
		return tint.NewHandler(output, &tint.Options{
			AddSource:   addSource,
			Level:       level,
			ReplaceAttr: replaceAttr,
			TimeFormat:  time.DateTime,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported log format %q", format)
	}
}

func resolveFormat(rawFormat, environment string) (string, error) {
	format := strings.ToLower(strings.TrimSpace(rawFormat))
	if format == "" || format == "auto" {
		if environment == "production" {
			return "json", nil
		}

		return "pretty", nil
	}

	switch format {
	case "pretty", "json":
		return format, nil
	default:
		return "", fmt.Errorf("invalid log format %q", rawFormat)
	}
}

func parseLevel(rawLevel string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(rawLevel)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid log level %q", rawLevel)
	}
}
