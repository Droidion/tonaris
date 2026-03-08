package logging

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Middleware(base *slog.Logger) func(http.Handler) http.Handler {
	if base == nil {
		base = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			writer := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			requestLogger := base.With(
				"component", "http",
				"request_id", chimiddleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", requestPath(r),
				"remote_ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			ctx := WithContext(r.Context(), requestLogger)
			r = r.WithContext(ctx)

			defer func() {
				if recovered := recover(); recovered != nil {
					requestLogger.ErrorContext(ctx, "panic recovered",
						"panic", fmt.Sprint(recovered),
						"stack_trace", string(debug.Stack()),
					)

					if !writer.wroteHeader {
						http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}

				duration := time.Since(start).Milliseconds()
				attrs := []any{
					"status", writer.status,
					"duration_ms", duration,
					"bytes", writer.bytes,
				}

				switch {
				case writer.status >= http.StatusInternalServerError:
					requestLogger.ErrorContext(ctx, "request completed", attrs...)
				case writer.status >= http.StatusBadRequest:
					requestLogger.WarnContext(ctx, "request completed", attrs...)
				default:
					requestLogger.InfoContext(ctx, "request completed", attrs...)
				}
			}()

			next.ServeHTTP(writer, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (writer *responseWriter) Unwrap() http.ResponseWriter {
	return writer.ResponseWriter
}

func (writer *responseWriter) WriteHeader(statusCode int) {
	if writer.wroteHeader {
		return
	}

	writer.status = statusCode
	writer.wroteHeader = true
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *responseWriter) Write(data []byte) (int, error) {
	if !writer.wroteHeader {
		writer.WriteHeader(http.StatusOK)
	}

	written, err := writer.ResponseWriter.Write(data)
	writer.bytes += written
	return written, err
}

func (writer *responseWriter) Flush() {
	flusher, ok := writer.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func (writer *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := writer.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support hijacking")
	}

	return hijacker.Hijack()
}

func (writer *responseWriter) Push(target string, options *http.PushOptions) error {
	pusher, ok := writer.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}

	return pusher.Push(target, options)
}

func requestPath(request *http.Request) string {
	if request.URL.Path == "" {
		return "/"
	}

	return request.URL.Path
}

var _ http.Flusher = (*responseWriter)(nil)
var _ http.Hijacker = (*responseWriter)(nil)
var _ http.Pusher = (*responseWriter)(nil)
var _ io.Writer = (*responseWriter)(nil)
