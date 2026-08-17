package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *task.TaskService
}

func NewTaskHandler(service task.TaskService) *TaskHandler {
	return &TaskHandler{service: &service}
}

func (h *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}

func (h *TaskHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	// TaskDefinition routes (nested under workflow_definition)
	r.Route("/workflow_definition/{workflowDefinitionID}/task_definition", func(r chi.Router) {
		r.Post("", h.createTaskDefinition)
		r.Get("", h.getTaskDefinitionsByWorkflowDefinition)
	})

	r.Route("/task_definition/{taskDefinitionID}", func(r chi.Router) {
		r.Get("", h.getTaskDefinitionByUserAndID)
		r.Put("", h.updateTaskDefinitionByID)
		r.Delete("", h.deleteTaskDefinition)
	})

	// TaskDependency routes
	r.Route("/task_dependency", func(r chi.Router) {
		r.Post("", h.createTaskDependency)
	})

	r.Route("/workflow_definition/{workflowDefinitionID}/task_dependency", func(r chi.Router) {
		r.Get("", h.getTaskDependencyByWorkflowDefinition)
	})

	r.Route("/task_dependency/{taskDependencyID}", func(r chi.Router) {
		r.Get("", h.getTaskDependencyByUserAndID)
		r.Put("", h.updateTaskDependencyByID)
		r.Delete("", h.deleteTaskDependency)
	})

	r.Route("/workflow_run/{workflowRunID}/task_run", func(r chi.Router) {
		r.Get("", h.getTaskRunByWorkflowRun)
	})

	r.Route("/task_run/{taskRunID}", func(r chi.Router) {
		r.Get("", h.getTaskRunByUserAndID)
		r.Put("", h.updateTaskRunByID)
		r.Delete("", h.deleteTaskRun)

		r.Route("/task_attempt", func(r chi.Router) {
			r.Get("", h.getTaskAttemptByTaskRun)
		})
	})

	// TaskAttempt routes
	r.Route("/task_attempt/{taskAttemptID}", func(r chi.Router) {
		r.Get("", h.getTaskAttemptByUserAndID)
		r.Delete("", h.deleteTaskAttempt)
	})

	return r
}

// TaskDefinition

func (h *TaskHandler) createTaskDefinition(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowDefinitionID, err := uuid.Parse(chi.URLParam(r, "workflowDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var tDRequestDto task.TaskDefinitionRequestDto
	err = render.DecodeJSON(r.Body, &tDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := h.service.CreateTaskDefinition(
		r.Context(),
		userID,
		workflowDefinitionID,
		tDRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (h *TaskHandler) getTaskDefinitionsByWorkflowDefinition(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowDefinitionID, err := uuid.Parse(chi.URLParam(r, "workflowDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := h.service.GetTaskDefinitionsByWorkflowDefinition(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (h *TaskHandler) getTaskDefinitionByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDefinitionID, err := uuid.Parse(chi.URLParam(r, "taskDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := h.service.GetTaskDefinitionByUserAndID(
		r.Context(),
		userID,
		taskDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (h *TaskHandler) updateTaskDefinitionByID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDefinitionID, err := uuid.Parse(chi.URLParam(r, "taskDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var tDRequestDto task.TaskDefinitionRequestDto
	err = render.DecodeJSON(r.Body, &tDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := h.service.UpdateTaskDefinitionByID(
		r.Context(),
		userID,
		taskDefinitionID,
		tDRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (h *TaskHandler) deleteTaskDefinition(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDefinitionID, err := uuid.Parse(chi.URLParam(r, "taskDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := h.service.SoftDeleteTaskDefinition(r.Context(), userID, taskDefinitionID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

// TaskDependency

func (h *TaskHandler) createTaskDependency(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	var tDpRequestDto task.TaskDependencyRequestDto
	err := render.DecodeJSON(r.Body, &tDpRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := h.service.CreateTaskDependency(r.Context(), userID, tDpRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (h *TaskHandler) getTaskDependencyByWorkflowDefinition(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowDefinitionID, err := uuid.Parse(chi.URLParam(r, "workflowDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := h.service.GetTaskDependencyByWorkflowDefinition(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (h *TaskHandler) getTaskDependencyByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDependencyID, err := uuid.Parse(chi.URLParam(r, "taskDependencyID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := h.service.GetTaskDependencyByUserAndID(
		r.Context(),
		userID,
		taskDependencyID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (h *TaskHandler) updateTaskDependencyByID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDependencyID, err := uuid.Parse(chi.URLParam(r, "taskDependencyID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var tDpRequestDto task.TaskDependencyRequestDto
	err = render.DecodeJSON(r.Body, &tDpRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := h.service.UpdateTaskDependencyByID(
		r.Context(),
		userID,
		taskDependencyID,
		tDpRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (h *TaskHandler) deleteTaskDependency(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskDependencyID, err := uuid.Parse(chi.URLParam(r, "taskDependencyID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := h.service.SoftDeleteTaskDependency(
		r.Context(),
		userID,
		taskDependencyID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

// TaskRun

func (h *TaskHandler) getTaskRunByWorkflowRun(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowRunID, err := uuid.Parse(chi.URLParam(r, "workflowRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := h.service.GetTaskRunByWorkflowRun(r.Context(), userID, workflowRunID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (h *TaskHandler) getTaskRunByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskRunID, err := uuid.Parse(chi.URLParam(r, "taskRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := h.service.GetTaskRunByUserAndID(r.Context(), userID, taskRunID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (h *TaskHandler) updateTaskRunByID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskRunID, err := uuid.Parse(chi.URLParam(r, "taskRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var tRRequestDto task.TaskRunRequestDto
	err = render.DecodeJSON(r.Body, &tRRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := h.service.UpdateTaskRunByUserAndID(
		r.Context(),
		userID,
		taskRunID,
		tRRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (h *TaskHandler) deleteTaskRun(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskRunID, err := uuid.Parse(chi.URLParam(r, "taskRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := h.service.SoftDeleteTaskRun(r.Context(), userID, taskRunID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

// TaskAttempt

func (h *TaskHandler) getTaskAttemptByTaskRun(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskRunID, err := uuid.Parse(chi.URLParam(r, "taskRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := h.service.GetTaskAttemptByTaskRun(r.Context(), userID, taskRunID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}

func (h *TaskHandler) getTaskAttemptByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskAttemptID, err := uuid.Parse(chi.URLParam(r, "taskAttemptID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := h.service.GetTaskAttemptByUserAndID(r.Context(), userID, taskAttemptID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}

func (h *TaskHandler) deleteTaskAttempt(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	taskAttemptID, err := uuid.Parse(chi.URLParam(r, "taskAttemptID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := h.service.SoftDeleteTaskAttempt(r.Context(), userID, taskAttemptID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}
