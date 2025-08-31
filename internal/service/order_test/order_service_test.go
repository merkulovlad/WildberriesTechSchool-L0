package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/merkulovlad/wbtech-go/internal/mocks"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"github.com/merkulovlad/wbtech-go/internal/service/order"
	"github.com/stretchr/testify/require"
)

func TestOrderService_GetOrder_CacheMissThenSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockInterfaceCache(ctrl)
	svc := order.NewOrderService(mockRepo, mockCache)

	id := "b1"
	expected := &model.Order{OrderUID: id, TrackNumber: "TRK001"}

	hit := false

	mockCache.EXPECT().
		Get(id).
		DoAndReturn(func(string) (*model.Order, bool) {
			if hit {
				return expected, true
			}
			return nil, false
		}).
		AnyTimes()

	mockRepo.EXPECT().
		GetOrder(gomock.Any(), id).
		Return(expected, nil).
		Times(1)

	mockCache.EXPECT().
		Set(gomock.Eq(id), gomock.AssignableToTypeOf(&model.Order{})).
		DoAndReturn(func(string, *model.Order) error {
			hit = true
			return nil
		}).
		Times(1)

	got, err := svc.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestOrderService_GetOrder_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockInterfaceCache(ctrl)

	svc := order.NewOrderService(mockRepo, mockCache)

	expected := &model.Order{OrderUID: "123", TrackNumber: "TRK001"}

	mockCache.EXPECT().
		Get("123").
		Return(expected, true)

	mockRepo.EXPECT().
		GetOrder(gomock.Any(), gomock.Any()).
		Times(0)

	got, err := svc.Get(context.Background(), "123")
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestOrderService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockInterfaceCache(ctrl)
	svc := order.NewOrderService(mockRepo, mockCache)

	ctx := context.Background()
	in := &model.Order{OrderUID: "o-1"}

	mockRepo.EXPECT().
		UpsertOrder(gomock.Any(), in).
		Return(nil).
		Times(1)

	err := svc.Create(ctx, in)
	require.NoError(t, err)
}

func TestOrderService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockInterfaceCache(ctrl)
	svc := order.NewOrderService(mockRepo, mockCache)

	ctx := context.Background()
	in := &model.Order{OrderUID: "o-1"}
	wantErr := errors.New("upsert failed")

	mockRepo.EXPECT().
		UpsertOrder(gomock.Any(), in).
		Return(wantErr).
		Times(1)

	err := svc.Create(ctx, in)
	require.Error(t, err)
	require.EqualError(t, err, wantErr.Error())
}

func TestOrderService_UpdateCache_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockCache := mocks.NewMockInterfaceCache(ctrl)
	svc := order.NewOrderService(mockRepo, mockCache)

	ctx := context.Background()
	o1 := &model.Order{OrderUID: "id-1"}
	o2 := &model.Order{OrderUID: "id-2"}
	recent := []*model.Order{o1, o2}

	mockRepo.EXPECT().
		GetRecent(gomock.Any(), 10).
		Return(recent, nil).
		Times(1)

	mockCache.EXPECT().
		Set("id-1", o1).
		Return(nil).
		Times(1)
	mockCache.EXPECT().
		Set("id-2", o2).
		Return(nil).
		Times(1)

	err := svc.UpdateCache(ctx)
	require.NoError(t, err)
}
