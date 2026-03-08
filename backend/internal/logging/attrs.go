package logging

import (
	"log/slog"
	"path/filepath"
)

func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == "" {
		return slog.Attr{}
	}

	if err, ok := attr.Value.Any().(error); ok {
		return slog.String(attr.Key, err.Error())
	}

	if attr.Key != slog.SourceKey {
		return attr
	}

	source := attrSource(attr.Value)
	if source == nil {
		return attr
	}

	source.File = shortenSourceFile(source.File)
	return slog.Any(attr.Key, source)
}

func attrSource(value slog.Value) *slog.Source {
	if value.Kind() != slog.KindAny {
		return nil
	}

	switch source := value.Any().(type) {
	case *slog.Source:
		if source == nil {
			return nil
		}

		copy := *source
		return &copy
	case slog.Source:
		copy := source
		return &copy
	default:
		return nil
	}
}

func shortenSourceFile(path string) string {
	if path == "" {
		return path
	}

	return filepath.Base(path)
}
