package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// HandlerFunc converts a handler with an error to a standard handler
func HandlerFunc(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log := GetLoggerFromContext(r.Context())
			log.Printf("Error from handler: %v", err)
		}
	})
}

// GetDefaultLogger gets a default logger to use
func GetDefaultLogger(serviceName, environment string) *logrus.Entry {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableTimestamp: false}
	return logrus.NewEntry(l).WithFields(logrus.Fields{
		"service":     serviceName,
		"environment": environment,
	})
}
