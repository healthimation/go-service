package middleware

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey int

// context keys
const (
	contextKeyLogger contextKey = iota
	contextKeyRequestID
	contextKeyInsecureUserID
)

// getStringFromContext returns a string from the context or empty string if missing
func getStringFromContext(ctx context.Context, key contextKey) string {
	stri := ctx.Value(key)
	if str, ok := stri.(string); ok {
		return str
	}
	return ""
}

// getRequestIDFromContext returns the requestID or empty string if missing
func getRequestIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextKeyRequestID)
}

// getInsecureUserIDFromContext returns the user id or empty string if missing
func getInsecureUserIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextKeyInsecureUserID)
}

// GetLoggerFromContext returns a logrus entry from the context, defaulting to a standard one if its not found
func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	inter := ctx.Value(contextKeyLogger)
	if entry, ok := inter.(*logrus.Entry); ok {
		return entry
	}
	return logrus.NewEntry(logrus.New()) // default to a new entry
}
