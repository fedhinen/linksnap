package storage

import (
	"database/sql"
	"fmt"
	"linksnap/internal/config"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(env *config.Env) (*sql.DB, error) {
	dbDriver := env.DatabaseDriver // postgres or sqlite3
	dbUrl := env.DatabaseURL

	db, err := sql.Open(dbDriver, dbUrl)

	if err != nil {
		return nil, fmt.Errorf("Error conectando con la base de datos %s", err)
	}

	return db, nil
}
