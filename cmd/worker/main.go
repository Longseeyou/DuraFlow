package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/adapters/kafka"
	"github.com/Longseeyou/DuraFlow/internal/adapters/postgresql"
	"github.com/Longseeyou/DuraFlow/internal/worker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	kafkaConfig := kafka.NewKafkaConfig(
		[]string{"localhost:9190", "localhost:9191", "localhost:9192"},
		3,
		time.Duration(10*time.Second),
		nil,
		nil,
	)
	c := kafka.NewKafkaConsumer("0", &kafkaConfig, "0", []string{"TaskCommandRequest"})
	p := kafka.NewKafkaProducer("0", &kafkaConfig)

	dsn := "host=localhost user=gorm password=gorm dbname=duraflow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	tRI := postgresql.NewPostgresTaskRepositoryInternal(db)

	ctx := context.Background()

	err = c.Start(ctx)
	if err != nil {
		slog.Error("")
		return
	}

	err = p.Start(ctx)
	if err != nil {
		return
	}

	w := worker.NewWorker("0", c, p, tRI)

	w.Run(ctx)
}
