package cache

import model "github.com/merkulovlad/wbtech-go/internal/model"

type InterfaceCache interface {
	Get(key string) (*model.Order, bool)
	Set(key string, value *model.Order) error
}
