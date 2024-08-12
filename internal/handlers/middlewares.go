package handlers

import (

	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

func GzipCompress(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        
		ow := w  // сохраняем оригинальный Writer
        
		// Проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
        acceptEncoding := r.Header.Get("Accept-Encoding")
        supportsGzip := strings.Contains(acceptEncoding, "gzip")

		// Проверяем подходщий тип контента
		contentType := r.Header.Get("Content-Type")
		_, found := ContentTypesForGzip[contentType]

        if supportsGzip && found {
            // оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
            cw := newCompressWriter(w)
            // меняем оригинальный http.ResponseWriter на новый
            ow = cw
			ow.Header().Set("Content-Encoding", "gzip")
            // не забываем отправить клиенту все сжатые данные после завершения middleware
            defer cw.Close()
        }

        // проверяем, что клиент отправил серверу сжатые данные в формате gzip
        contentEncoding := r.Header.Get("Content-Encoding")
        sendsGzip := strings.Contains(contentEncoding, "gzip")
        if sendsGzip {
            // оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
            cr, err := newCompressReader(r.Body)
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
            // меняем тело запроса на новое
            r.Body = cr
            defer cr.Close()
        }

        // передаём управление хендлеру
        h.ServeHTTP(ow, r)
    }
}
