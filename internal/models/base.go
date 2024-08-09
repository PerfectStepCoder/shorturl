package models

type RequestFullURL struct {
	URL string `json:"url"`
}

type ResponseShortURL struct {
	Result string `json:"result"`
}
