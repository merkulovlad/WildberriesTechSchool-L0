package order

import (
	"context"

	"github.com/merkulovlad/wbtech-go/internal/model"
)

type Service interface {
	Get(c context.Context, id string) (*model.Order, error)
	UpdateCache(c context.Context) error
	Create(c context.Context, order *model.Order) error
}
