package middleware

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	logFieldRequestID = "request_id"
	logFieldUserID    = "user_id"
)

// Logger injects a logger into the context
type Logger interface {
	Log(h http.Handler) http.Handler
}

type logger struct {
	entry       *logrus.Entry
	logRequests bool
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.StatusCode = code
	l.ResponseWriter.WriteHeader(code)
}

// NewLogrusLogger allows you to setup a base entry to use for logging
func NewLogrusLogger(baseEntry *logrus.Entry, logRequests bool) Logger {
	return &logger{entry: baseEntry, logRequests: logRequests}
}

func (l *logger) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := l.entry.WithFields(logrus.Fields{
			logFieldRequestID: getRequestIDFromContext(r.Context()),
			logFieldUserID:    getInsecureUserIDFromContext(r.Context()),
		})
		// add logger to the context
		ctx := context.WithValue(r.Context(), contextKeyLogger, entry)
		r = r.WithContext(ctx)

		lrw := &loggingResponseWriter{w, http.StatusOK}
		h.ServeHTTP(lrw, r)
		if l.logRequests {
			entry.Printf("%s - [%d] - %s", r.Method, lrw.StatusCode, r.URL.String())
		}
	})
}
