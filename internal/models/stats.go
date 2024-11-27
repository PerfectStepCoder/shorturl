// Модуль models содержит описание получаемых и возвращаемых сущностей HTTP сервисом.
package models

// ResponseStatsBase - возвращаемый результат обработки статситики обработанных ссылок пользователями.
type ResponseStatsBase struct {
	Urls  int `json:"urls"`  // количество сокращённых URL в сервисе
	Users int `json:"users"` // количество пользователей в сервисе
}
