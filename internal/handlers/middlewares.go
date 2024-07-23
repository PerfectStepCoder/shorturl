package handlers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func WithLogging(h http.HandlerFunc, logger *logrus.Logger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		duration := time.Since(start)

		logData := map[string]interface{}{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}

		logger.WithFields(logrus.Fields(logData)).Info()
	}
	// Возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
