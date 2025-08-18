// repository/connection.go
package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/merkulovlad/wbtech-go/internal/config/config"
)

const Driver = "postgres"

func ConnectDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.DSN()
	db, err := sql.Open(Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return db, nil
}
