package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/imabg/onboard/internal/logger"
	"go.uber.org/zap"
)

const maxPayloadBytes = 1 << 20

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		payload := readPayload(r)

		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.L().Info("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("route", r.URL.RequestURI()),
			zap.String("host", r.Host),
			zap.String("proto", r.Proto),
			zap.String("remote", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("content_type", r.Header.Get("Content-Type")),
			zap.Int64("content_length", r.ContentLength),
			payloadField(payload),
			zap.Int("status", rw.status),
			zap.Int("bytes", rw.bytes),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func readPayload(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxPayloadBytes+1))
	_ = r.Body.Close()
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(nil))
		return nil
	}
	if len(body) > maxPayloadBytes {
		body = body[:maxPayloadBytes]
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body
}

func payloadField(body []byte) zap.Field {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return zap.String("payload", "")
	}
	if json.Valid(trimmed) {
		var v any
		if err := json.Unmarshal(trimmed, &v); err == nil {
			return zap.Any("payload", v)
		}
	}
	return zap.String("payload", string(trimmed))
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
