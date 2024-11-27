// Модуль models содержит описание получаемых и возвращаемых сущностей HTTP сервисом.
package models

type ResponseStatsBase struct {
	Urls  int `json:"urls"`  // количество сокращённых URL в сервисе
	Users int `json:"users"` // количество пользователей в сервисе
}
