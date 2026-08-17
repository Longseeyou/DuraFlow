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
	producer := kafka.NewKafkaProducer(
		envString("WORKER_PRODUCER_ID", "worker-0"),
		&kafkaConfig,
	)

	db, err := gorm.Open(
		postgres.Open(envString("DATABASE_DSN", "host=localhost user=gorm password=gorm dbname=duraflow")),
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

	if err := producer.Start(ctx); err != nil {
		slog.Error("start producer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := producer.Stop(); err != nil {
			slog.Error("stop producer", "error", err)
		}
	}()

	taskRepositoryInternal := postgresql.NewPostgresTaskRepositoryInternal(db)

	w := worker.NewWorker(
		envString("WORKER_ID", "worker-0"),
		consumer,
		producer,
		taskRepositoryInternal,
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
