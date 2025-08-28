package repository

import (
	"context"

	"github.com/merkulovlad/wbtech-go/internal/model"
)

type Repository interface {
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	GetRecent(ctx context.Context, limit int) ([]*model.Order, error)
	UpsertOrder(ctx context.Context, o *model.Order) error
}
