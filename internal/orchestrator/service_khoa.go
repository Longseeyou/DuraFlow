package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type OrchestratorService struct {
	taskRepository task.TaskRepositoryInternal
	taskPoller     task.TaskPoller

	consumer message.Consumer
	producer message.Producer
}

func NewOrchestratorService(
	taskRepository task.TaskRepositoryInternal,
	taskPoller task.TaskPoller,
	consumer message.Consumer,
	producer message.Producer,
) *OrchestratorService {
	return &OrchestratorService{
		taskRepository: taskRepository,
		taskPoller:     taskPoller,
		consumer:       consumer,
		producer:       producer,
	}
}

func (oS *OrchestratorService) Run(ctx context.Context) {
}

func (oS *OrchestratorService) SendTaskCommandRequest(ctx context.Context) {
	for {
		// NOTE: Poll task
		tR, err := oS.taskPoller.PollTask(ctx)
		if err != nil {
			// TODO: err
		}

		// NOTE: Idempotency
		valid, err := oS.taskRepository.TaskRunIdempotency(ctx, tR.ID, task.TASK_RUN_QUEUED)
		if err != nil {
			// TODO: err
		}
		if !valid {
			// TODO:
		}

		// NOTE: Create TaskAttempt
		var tA task.TaskAttempt
		tA.TaskRunID = tR.ID

		tA, err = oS.taskRepository.CreateTaskAttempt(ctx, tA)
		if err != nil {
			// TODO:
		}

		// NOTE: Send TaskCommandRequest
		var tCRequest task.TaskCommandRequest
		tCRequest.TaskAttemptID = tA.ID
		tCRequest.TaskType = tR.TaskDefinition.TaskType
		tCRequest.Input = tR.Input
		tCRequest.Timeout = &tR.TaskDefinition.Timeout

		value, err := json.Marshal(tCRequest)
		if err != nil {
			// TODO:
		}

		err = oS.producer.SendMessage(
			ctx,
			"TaskCommandRequest",
			"TaskCommandRequest",
			value,
		)
		if err != nil {
			// TODO:
		}
	}
}

func (oS *OrchestratorService) ReceiveTaskCommandResponse(ctx context.Context) {
	for {
		message, err := oS.consumer.ReceiveMessage(ctx)
		if err != nil {
			// TODO:
		}

		switch string(message.Key) {
		case "TaskCommandResponse":
			var tCResponse task.TaskCommandResponse
			err := json.Unmarshal(message.Value, &tCResponse)
			if err != nil {
				// TODO:
			}

			valid, err := oS.taskRepository.TaskAttemptIdempotency(
				ctx,
				tCResponse.TaskAttemptID,
				tCResponse.Status,
			)
			if err != nil {
				// TODO:
			}
			if !valid {
				// TODO:
			}

			oS.HandleTaskCommandResponse(ctx, tCResponse)
		}
	}
}

func (oS *OrchestratorService) HandleTaskCommandResponse(
	ctx context.Context,
	tcResponse task.TaskCommandResponse,
) {
	newTA := make(map[string]any)
	newTA["status"] = tcResponse.Status
	newTA["worker_id"] = tcResponse.WorkerID
	newTA["started_at"] = tcResponse.StartedAt
	newTA["ended_at"] = tcResponse.EndedAt
	if tcResponse.Log != nil {
		newTA["log"] = *tcResponse.Log
	}

	tA, err := oS.taskRepository.UpdateTaskAttemptById(ctx, tcResponse.TaskAttemptID, newTA)
	if err != nil {
		// TODO:
	}

	newTR := make(map[string]any)

	switch tcResponse.Status {
	case task.TASK_ATTEMPT_COMPLETED:
		newTR["status"] = task.TASK_RUN_COMPLETED
		newTR["ended_at"] = *tA.EndedAt
		if tcResponse.Output != nil {
			newTR["output"] = *tcResponse.Output
		}

	case task.TASK_ATTEMPT_FAILED:
		newTR["status"] = task.TASK_RUN_FAILED
	default:
	}

	_, err = oS.taskRepository.UpdateTaskRunById(ctx, tA.TaskRunID, newTR)
	if err != nil {
		// TODO:
	}
}
