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
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultTaskRunCDCTopic = "dbz.public.task_runs"

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
		envString("ORCHESTRATOR_ID", "orchestrator-0"),
		&kafkaConfig,
		envString("ORCHESTRATOR_GROUP_ID", "orchestrator"),
		[]string{envString("CDC_TASK_RUN_TOPIC", defaultTaskRunCDCTopic)},
	)
	producer := kafka.NewKafkaProducer(
		envString("ORCHESTRATOR_PRODUCER_ID", "orchestrator-0"),
		&kafkaConfig,
	)

	db, err := gorm.Open(
		postgres.Open(
			envString("DATABASE_DSN", "host=localhost user=gorm password=gorm dbname=duraflow"),
		),
		&gorm.Config{},
	)
	if err != nil {
		slog.Error("connect database", "error", err)
		os.Exit(1)
	}
	if err := migrate(db); err != nil {
		slog.Error("migrate database", "error", err)
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

	repository := postgresql.NewPostgresOrchestratorRepository(db)
	taskPoller := kafka.NewKafkaTaskPoller(consumer)
	rescueTaskPoller := postgresql.NewPostgresTaskPoller(db)

	scheduler := orchestrator.NewTaskScheduler(repository, taskPoller, rescueTaskPoller, producer)
	orch := orchestrator.NewOrchestrator(repository)

	slog.Info(
		"orchestrator started",
		"task_command_topic",
		orchestrator.TaskCommandRequestTopic,
		"task_command_response_topic",
		orchestrator.TaskCommandResponseTopic,
	)

	go scheduler.Run(ctx)
	go orch.ResolveTimeoutTaskRun(ctx)

	<-ctx.Done()
	slog.Info("orchestrator stopped")
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&workflow.Workflow{},
		&workflow.WorkflowDefinition{},
		&workflow.WorkflowRun{},
		&task.TaskDefinition{},
		&task.TaskDependency{},
		&task.TaskRun{},
		&task.TaskAttempt{},
		// &task.TaskEvent{},
	)
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
