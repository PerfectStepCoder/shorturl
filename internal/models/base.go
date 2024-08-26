package models

type RequestFullURL struct {
	URL string `json:"url"`
}

type ResponseShortURL struct {
	Result string `json:"result"`
}

type RequestCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseCorrelationURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
