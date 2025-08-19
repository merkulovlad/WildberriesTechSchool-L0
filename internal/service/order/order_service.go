package order

import (
	"context"
	"github.com/merkulovlad/wbtech-go/internal/db/repository"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"github.com/merkulovlad/wbtech-go/internal/service/cache"
	"golang.org/x/sync/singleflight"
)

type orderService struct {
	repo  repository.Repository
	cache cache.InterfaceCache
	group singleflight.Group
}

func NewOrderService(r repository.Repository, c cache.InterfaceCache) Service {
	return &orderService{
		repo:  r,
		cache: c,
		group: singleflight.Group{},
	}
}

func (s *orderService) Get(c context.Context, id string) (*model.Order, error) {
	if order, exists := s.cache.Get(id); exists {
		return order, nil
	}
	res, err, _ := s.group.Do(id, func() (interface{}, error) {
		if order, exists := s.cache.Get(id); exists {
			return order, nil
		}

		order, err := s.repo.GetOrder(c, id)
		if err != nil {
			return nil, err
		}

		_ = s.cache.Set(id, order)
		return order, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*model.Order), nil
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
