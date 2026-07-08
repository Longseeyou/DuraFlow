package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

const (
	DefaultWorkflowRunTopic = "WorkflowRunRequest"
	DefaultTaskEventTopic   = "TaskEvent"
)

type SchedulerOption func(*Scheduler)

func WithSchedulerTopics(workflowRunTopics []string, taskEventTopics []string) SchedulerOption {
	return func(scheduler *Scheduler) {
		scheduler.workflowRunTopics = routeSet(workflowRunTopics)
		scheduler.taskEventTopics = routeSet(taskEventTopics)
	}
}

func WithSchedulerLogger(logger *slog.Logger) SchedulerOption {
	return func(scheduler *Scheduler) {
		if logger != nil {
			scheduler.logger = logger
		}
	}
}

type Scheduler struct {
	service           *Service
	consumer          message.Consumer
	workflowRunTopics map[string]struct{}
	taskEventTopics   map[string]struct{}
	logger            *slog.Logger
}

func NewScheduler(
	service *Service,
	consumer message.Consumer,
	options ...SchedulerOption,
) *Scheduler {
	scheduler := &Scheduler{
		service:           service,
		consumer:          consumer,
		workflowRunTopics: routeSet([]string{DefaultWorkflowRunTopic}),
		taskEventTopics:   routeSet([]string{DefaultTaskEventTopic}),
		logger:            slog.Default(),
	}
	for _, option := range options {
		option(scheduler)
	}
	return scheduler
}

func (scheduler *Scheduler) Start(ctx context.Context) error {
	if err := scheduler.consumer.Start(ctx); err != nil {
		return fmt.Errorf("start scheduler consumer: %w", err)
	}
	defer func() {
		if err := scheduler.consumer.Stop(); err != nil {
			scheduler.logger.Error("stop scheduler consumer", "error", err)
		}
	}()
	return scheduler.Run(ctx)
}

func (scheduler *Scheduler) Run(ctx context.Context) error {
	for {
		msg, err := scheduler.consumer.ReceiveMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			scheduler.logger.Error("receive scheduler message", "error", err)
			continue
		}
		if msg == nil {
			continue
		}

		if err := scheduler.ProcessMessage(ctx, *msg); err != nil {
			scheduler.logger.Error(
				"process scheduler message",
				"topic",
				msg.Topic,
				"key",
				string(msg.Key),
				"partition",
				msg.Partition,
				"offset",
				msg.Offset,
				"error",
				err,
			)
		}
	}
}

func (scheduler *Scheduler) ProcessMessage(ctx context.Context, msg message.Message) error {
	if scheduler.service == nil {
		return errors.New("scheduler service is required")
	}

	route := messageRoute(msg)
	if routeMatches(route, scheduler.taskEventTopics) {
		event, _, err := decodeTaskEvent(msg, true)
		if err != nil {
			return err
		}
		return scheduler.service.HandleTaskEvent(ctx, event)
	}
	if routeMatches(route, scheduler.workflowRunTopics) {
		workflowRunID, err := decodeWorkflowRunID(msg)
		if err != nil {
			return err
		}
		return scheduler.service.ScheduleWorkflow(ctx, workflowRunID)
	}

	event, ok, err := decodeTaskEvent(msg, false)
	if err != nil {
		return err
	}
	if ok {
		return scheduler.service.HandleTaskEvent(ctx, event)
	}

	workflowRunID, err := decodeWorkflowRunID(msg)
	if err != nil {
		return err
	}
	return scheduler.service.ScheduleWorkflow(ctx, workflowRunID)
}

func decodeTaskEvent(msg message.Message, required bool) (task.TaskEvent, bool, error) {
	var event task.TaskEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		if required {
			return task.TaskEvent{}, false, fmt.Errorf("decode task event: %w", err)
		}
		return task.TaskEvent{}, false, nil
	}
	if event.EventType == "" {
		if required {
			return task.TaskEvent{}, false, errors.New("task event type is required")
		}
		return task.TaskEvent{}, false, nil
	}
	return event, true, nil
}

func decodeWorkflowRunID(msg message.Message) (uuid.UUID, error) {
	trimmed := bytes.TrimSpace(msg.Value)
	if len(trimmed) > 0 {
		if id, err := uuid.Parse(strings.Trim(string(trimmed), `"`)); err == nil {
			return id, nil
		}

		var envelope struct {
			WorkflowRunID uuid.UUID `json:"workflow_run_id"`
			ID            uuid.UUID `json:"id"`
		}
		if err := json.Unmarshal(trimmed, &envelope); err != nil {
			return uuid.Nil, fmt.Errorf("decode workflow run request: %w", err)
		}
		if envelope.WorkflowRunID != uuid.Nil {
			return envelope.WorkflowRunID, nil
		}
		if envelope.ID != uuid.Nil {
			return envelope.ID, nil
		}
	}

	if id, err := uuid.Parse(string(msg.Key)); err == nil {
		return id, nil
	}
	return uuid.Nil, errors.New("workflow run ID is required")
}

func messageRoute(msg message.Message) string {
	if msg.Topic != "" {
		return msg.Topic
	}
	return string(msg.Key)
}

func routeSet(routes []string) map[string]struct{} {
	result := make(map[string]struct{}, len(routes))
	for _, route := range routes {
		normalized := normalizeRoute(route)
		if normalized != "" {
			result[normalized] = struct{}{}
		}
	}
	return result
}

func routeMatches(route string, routes map[string]struct{}) bool {
	_, ok := routes[normalizeRoute(route)]
	return ok
}

func normalizeRoute(route string) string {
	return strings.ToLower(strings.TrimSpace(route))
}
