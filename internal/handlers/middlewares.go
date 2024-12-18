// Модуль содержит декораторы для обработчиков.
package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

var (
	// TODO вынести данные переменные в переменные окружения
	hashKey  = []byte("very-secret-key")  // Симметричный ключ для подписи
	blockKey = []byte("a-lot-secret-key") // Симметричный ключ для шифрования
	sCookie  = securecookie.New(hashKey, blockKey)
)

type contextKey string

// UserKeyUID - идентификатор пользователя который передается в контексте.
const UserKeyUID contextKey = "userUID"

// WithLogging - декоратор логирование пользователя.
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

// GzipCompress - декоратор сжатия данных.
func GzipCompress(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ow := w // сохраняем оригинальный Writer

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

// ValidateUserUID - декоратор валидации JWT токена.
func ValidateUserUID(cookieValue string) (string, bool) {
	var userUID string
	if err := sCookie.Decode("userUID", cookieValue, &userUID); err != nil {
		return "", false
	}
	// Кука существует и проходит проверку, продолжаем выполнение следующего обработчика
	logrus.Printf("Existing valid user ID: %s", userUID)
	return userUID, true
}

// CheckSignedCookie для подписанной куки с идентификатором пользователя.
func CheckSignedCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Попытка получить куку
		cookie, err := r.Cookie("userUID")

		if err != nil {
			h.ServeHTTP(w, r)
		} else {
			// Проверка и декодирование куки
			userUID, isValid := ValidateUserUID(cookie.Value)

			if isValid {
				// Кука существует и проходит проверку, продолжаем выполнение следующего обработчика
				logrus.Printf("Existing valid user ID: %s", userUID)
				h.ServeHTTP(w, r)
				return
			} else {
				logrus.Printf("Wrong UserUID: %s", cookie.Value)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
}

// SetNewCookie - декоратор установки новой куки.
func SetNewCookie(w http.ResponseWriter) (string, error) {
	// Если куки нет или она невалидна, создаем новую
	userUID := uuid.New().String()

	// Кодирование и подпись куки
	encoded, err := sCookie.Encode("userUID", userUID)
	if err != nil {
		http.Error(w, "Error signing the cookie", http.StatusInternalServerError)
		return "", err
	}

	// Установка куки
	http.SetCookie(w, &http.Cookie{
		Name:  "userUID",
		Value: encoded,
		Path:  "/",
		// Опциональные параметры безопасности:
		HttpOnly: true,  // Доступ только через HTTP
		Secure:   false, // Отправка только по HTTPS
	})

	w.Header().Set("Authorization", encoded)

	logrus.Println("New user UID assigned:", userUID)
	return userUID, nil
}

// Auth для подписанной куки с идентификатором пользователя.
func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Попытка получить куку
		cookie, err := r.Cookie("userUID")

		if err != nil {
			if r.Method == http.MethodGet && r.URL.Path == "/api/user/urls" {
				encodedUserUID := r.Header.Get("Authorization")
				var validErr bool
				userUID, validErr := ValidateUserUID(encodedUserUID)
				if validErr {
					ctx := context.WithValue(r.Context(), UserKeyUID, userUID)
					h.ServeHTTP(w, r.WithContext(ctx))
				} else {
					logrus.Printf("Wrong UserUID: %s", encodedUserUID)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			} else { // Создаем пользователю uid (методы POST DELETE PUT)
				userUID, err := SetNewCookie(w)
				if err != nil {
					return
				}
				ctx := context.WithValue(r.Context(), UserKeyUID, userUID)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			// Проверка и декодирование куки
			userUID, isValid := ValidateUserUID(cookie.Value)
			if isValid {
				// Кука существует и проходит проверку, продолжаем выполнение следующего обработчика
				logrus.Printf("Existing valid user ID: %s", userUID)
				ctx := context.WithValue(r.Context(), UserKeyUID, userUID)
				h.ServeHTTP(w, r.WithContext(ctx))
			} else {
				logrus.Printf("Wrong UserUID: %s", cookie.Value)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
}
