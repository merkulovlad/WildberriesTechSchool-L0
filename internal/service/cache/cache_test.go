package cache

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/merkulovlad/wbtech-go/internal/mocks"
	"github.com/merkulovlad/wbtech-go/internal/model"
)

func TestCache_Get_Miss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := mocks.NewMockInterfaceLogger(ctrl)
	mockLog.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	c := NewCache(mockLog)

	got, ok := c.Get("missing")
	if ok {
		t.Fatalf("expected miss")
	}
	if got != nil {
		t.Fatalf("expected nil on miss, got %#v", got)
	}
}

func TestCache_Set_Then_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := mocks.NewMockInterfaceLogger(ctrl)
	mockLog.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	c := NewCache(mockLog)

	want := &model.Order{OrderUID: "k1", TrackNumber: "TRK001"}
	if err := c.Set("k1", want); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	got, ok := c.Get("k1")
	if !ok {
		t.Fatalf("expected hit")
	}
	if got != want {
		t.Fatalf("got %p, want %p", got, want)
	}
	if len(c.data) != 1 || c.order.Len() != 1 {
		t.Fatalf("sizes: data=%d order=%d", len(c.data), c.order.Len())
	}
}

func TestCache_UpdateExisting_DoesNotGrow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := mocks.NewMockInterfaceLogger(ctrl)
	mockLog.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	c := NewCache(mockLog)

	v1 := &model.Order{OrderUID: "k1", TrackNumber: "TRK001"}
	v2 := &model.Order{OrderUID: "k1", TrackNumber: "TRK001-UPDATED"}

	if err := c.Set("k1", v1); err != nil {
		t.Fatalf("Set v1: %v", err)
	}
	if len(c.data) != 1 || c.order.Len() != 1 {
		t.Fatalf("after v1 sizes: data=%d order=%d", len(c.data), c.order.Len())
	}

	if err := c.Set("k1", v2); err != nil {
		t.Fatalf("Set v2: %v", err)
	}

	if len(c.data) != 1 || c.order.Len() != 1 {
		t.Fatalf("after update sizes: data=%d order=%d", len(c.data), c.order.Len())
	}
	got, ok := c.Get("k1")
	if !ok || got != v2 {
		t.Fatalf("updated value not returned")
	}
}

func TestCache_Eviction_FIFO(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := mocks.NewMockInterfaceLogger(ctrl)
	mockLog.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	c := NewCache(mockLog)
	c.limit = 2

	a := &model.Order{OrderUID: "A"}
	b := &model.Order{OrderUID: "B"}
	c1 := &model.Order{OrderUID: "C"}

	if err := c.Set("A", a); err != nil {
		t.Fatalf("Set A: %v", err)
	}
	if err := c.Set("B", b); err != nil {
		t.Fatalf("Set B: %v", err)
	}
	if len(c.data) != 2 || c.order.Len() != 2 {
		t.Fatalf("fill sizes: data=%d order=%d", len(c.data), c.order.Len())
	}

	if err := c.Set("C", c1); err != nil {
		t.Fatalf("Set C: %v", err)
	}

	if _, ok := c.Get("A"); ok {
		t.Fatalf("A should be evicted (FIFO)")
	}
	if got, ok := c.Get("B"); !ok || got != b {
		t.Fatalf("expected B present")
	}
	if got, ok := c.Get("C"); !ok || got != c1 {
		t.Fatalf("expected C present")
	}
	if len(c.data) != 2 || c.order.Len() != 2 {
		t.Fatalf("post-evict sizes: data=%d order=%d", len(c.data), c.order.Len())
	}
}

func TestCache_Concurrent_SetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLog := mocks.NewMockInterfaceLogger(ctrl)
	mockLog.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	c := NewCache(mockLog)
	c.limit = 100

	keys := []string{"x1", "x2", "x3", "x4", "x5"}
	var wg sync.WaitGroup

	for _, k := range keys {
		wg.Add(1)
		k := k
		go func() {
			defer wg.Done()
			_ = c.Set(k, &model.Order{OrderUID: k})
		}()
	}
	for range keys {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = c.Get("x3")
			_, _ = c.Get("x-nope")
		}()
	}
	wg.Wait()

	if _, ok := c.Get("x3"); !ok {
		t.Fatalf("x3 should be present")
	}
}
