package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "github.com/merkulovlad/wbtech-go/internal/config/config"
	"github.com/merkulovlad/wbtech-go/internal/db/repository"
	"github.com/merkulovlad/wbtech-go/internal/kafka"
	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/server"
	"github.com/merkulovlad/wbtech-go/internal/service/cache"
	"github.com/merkulovlad/wbtech-go/internal/service/order"

	"github.com/gofiber/swagger"
	_ "github.com/merkulovlad/wbtech-go/docs"
)

// @title           Order Service API
// @version         1.0
// @description     API for managing and retrieving orders.
// @host            localhost:8080
// @BasePath        /
// @schemes         http
func main() {
	config := cfg.MustLoad()
	log, err := logger.NewLogger(&config.Log)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer func(log *logger.Logger) {
		if err := log.Sync(); err != nil {
			fmt.Printf("failed to sync logger: %v", err)
		}
	}(log)

	log.Info("loading database ")
	db, err := repository.ConnectDB(&config.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("failed to close database connection: %v", err)
		}
	}(db)
	log.Info("Migrating database ")
	err = repository.RunMigrations(db)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db, log)
	c := cache.NewCache(log)

	orderService := order.NewOrderService(orderRepo, c)
	consumer := kafka.NewConsumer(config.Kafka.Brokers, config.Kafka.Topic, config.Kafka.Group, "kafka.DLQ", orderService, log)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := consumer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Errorf("consumer stopped: %v", err)
	}
	ctxUpdate, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = orderService.UpdateCache(ctxUpdate); err != nil {
		log.Errorf("failed to update cache: %v", err)
	}

	log.Info("starting server")
	app := server.NewServer(orderService, log)
	app.Get("/swagger/*", swagger.HandlerDefault)
	go func() {
		if err := app.Listen(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
