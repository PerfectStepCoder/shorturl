// Модуль models содержит описание получаемых и возвращаемых сущностей HTTP сервисом.
package models

// RequestFullURL - передача полной ссылке для обработки.
type RequestFullURL struct {
	URL string `json:"url"`
}

// ResponseShortURL - возвращаемая короткая ссылка.
type ResponseShortURL struct {
	Result string `json:"result"`
}

// RequestCorrelationURL - запрос на обработку полной ссылке с идентификатором.
type RequestCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResponseCorrelationURL - возвращаемый результат обработки полной ссылке с идентификатором.
type ResponseCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ResponseURL - полный ответ, вклучая оригинальную ссылку и короткую.
type ResponseURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
