package order

import (
	"context"
	"github.com/merkulovlad/wbtech-go/internal/db/repository"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"github.com/merkulovlad/wbtech-go/internal/service/cache"
)

type orderService struct {
	repo  repository.Repository
	cache cache.InterfaceCache
}

func NewOrderService(r repository.Repository, c cache.InterfaceCache) Service {
	return &orderService{
		repo:  r,
		cache: c,
	}
}

func (s *orderService) Get(c context.Context, id string) (*model.Order, error) {
	order, exists := s.cache.Get(id)
	if exists {
		return order, nil
	}
	order, err := s.repo.GetOrder(c, id)
	if err != nil {
		return nil, err
	}
	err = s.cache.Set(id, order)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (s *orderService) Create(c context.Context, order *model.Order) error {
	return s.repo.UpsertOrder(c, order)
}

func (s *orderService) UpdateCache(c context.Context) error {
	orders, err := s.repo.GetRecent(c, 10)
	if err != nil {
		return err
	}
	for _, order := range orders {
		err = s.cache.Set(order.OrderUID, order)
		if err != nil {
			return err
		}
	}
	return nil
}
