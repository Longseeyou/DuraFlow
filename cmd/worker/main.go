package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/adapters/kafka"
	"github.com/Longseeyou/DuraFlow/internal/adapters/postgresql"
	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"github.com/Longseeyou/DuraFlow/internal/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	kafkaConfig := kafka.NewKafkaConfig(
		envList("KAFKA_BROKERS", []string{"localhost:9190", "localhost:9191", "localhost:9192"}),
		3,
		10*time.Second,
		envPointer("KAFKA_SASL_USER"),
		envPointer("KAFKA_SASL_PASS"),
	)

	consumer := kafka.NewKafkaConsumer(
		envString("WORKER_ID", "worker-0"),
		&kafkaConfig,
		envString("WORKER_GROUP_ID", "workers"),
		[]string{orchestrator.TaskCommandRequestTopic},
	)

	db, err := gorm.Open(
		postgres.Open(
			envString(
				"DATABASE_DSN",
				"host=localhost user=duraflow password=duraflow dbname=duraflow",
			),
		),
		&gorm.Config{},
	)
	if err != nil {
		slog.Error("connect database", "error", err)
		os.Exit(1)
	}

	if err := consumer.Start(ctx); err != nil {
		slog.Error("start consumer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := consumer.Stop(); err != nil {
			slog.Error("stop consumer", "error", err)
		}
	}()

	w := worker.NewWorker(
		envString("WORKER_ID", "worker-0"),
		consumer,
		postgresql.NewPostgresqlWorkerRepository(db),
	)

	w.Run(ctx)

	slog.Info("worker stopped")
}

func envString(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envList(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}

func envPointer(key string) *string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}
	return &value
}
