package repository

import (
	"database/sql"

	"github.com/merkulovlad/wbtech-go/internal/logger"
)

type OrderRepository struct {
	db     *sql.DB
	logger logger.InterfaceLogger
}

var _ Repository = (*OrderRepository)(nil)

func NewOrderRepository(db *sql.DB, log logger.InterfaceLogger) Repository {
	return &OrderRepository{
		db:     db,
		logger: log,
	}
}
