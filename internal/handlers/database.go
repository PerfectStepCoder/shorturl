package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// PingDatabase - обработчик проверки доступности базы данных.
func PingDatabase(databaseDSN string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		conn, err := pgx.Connect(context.Background(), databaseDSN)

		if err != nil {
			log.Print(err)
			http.Error(res, "Connect to DB not work", http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)

		defer conn.Close(context.Background())
	}
}
