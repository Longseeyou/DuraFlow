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

func (tH *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}

func (tH *TaskHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	// TaskDefinition routes (nested under workflow_definition)
	r.Route("/workflow_definition/{wDId}/task_definition", func(r chi.Router) {
		r.Post("", tH.createTaskDefinition)
		r.Get("", tH.getTaskDefinitionByWorkflowDefinition)
	})

	r.Route("/task_definition/{tDId}", func(r chi.Router) {
		r.Get("", tH.getTaskDefinitionByUserAndId)
		r.Put("", tH.updateTaskDefinitionById)
		r.Delete("", tH.deleteTaskDefinition)
	})

	// TaskDependency routes
	r.Route("/task_dependency", func(r chi.Router) {
		r.Post("", tH.createTaskDependency)
	})

	r.Route("/workflow_definition/{wDId}/task_dependency", func(r chi.Router) {
		r.Get("", tH.getTaskDependencyByWorkflowDefinition)
	})

	r.Route("/task_dependency/{tDpId}", func(r chi.Router) {
		r.Get("", tH.getTaskDependencyByUserAndId)
		r.Put("", tH.updateTaskDependencyById)
		r.Delete("", tH.deleteTaskDependency)
	})

	// TaskRun routes (nested under workflow_run)
	r.Route("/workflow_run/{wRId}/task_definition/{tDId}/task_run", func(r chi.Router) {
		r.Post("", tH.createTaskRun)
	})

	r.Route("/workflow_run/{wRId}/task_run", func(r chi.Router) {
		r.Get("", tH.getTaskRunByWorkflowRun)
	})

	r.Route("/task_run/{tRId}", func(r chi.Router) {
		r.Get("", tH.getTaskRunByUserAndId)
		r.Put("", tH.updateTaskRunById)
		r.Delete("", tH.deleteTaskRun)

		r.Route("/task_attempt", func(r chi.Router) {
			r.Get("", tH.getTaskAttemptByTaskRun)
		})
	})

	// TaskAttempt routes
	r.Route("/task_attempt/{tAId}", func(r chi.Router) {
		r.Get("", tH.getTaskAttemptByUserAndId)
		r.Delete("", tH.deleteTaskAttempt)
	})

	return r
}

// TaskDefinition

func (tH *TaskHandler) createTaskDefinition(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wDId, err := uuid.Parse(chi.URLParam(r, "wDId"))
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

	tDResponseDto, err := tH.service.CreateTaskDefinition(r.Context(), uId, wDId, tDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (tH *TaskHandler) getTaskDefinitionByWorkflowDefinition(
	w http.ResponseWriter,
	r *http.Request,
) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wDId, err := uuid.Parse(chi.URLParam(r, "wDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := tH.service.GetTaskDefinitionByWorkflowDefinition(r.Context(), uId, wDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (tH *TaskHandler) getTaskDefinitionByUserAndId(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDId, err := uuid.Parse(chi.URLParam(r, "tDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := tH.service.GetTaskDefinitionByUserAndId(r.Context(), uId, tDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (tH *TaskHandler) updateTaskDefinitionById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDId, err := uuid.Parse(chi.URLParam(r, "tDId"))
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

	tDResponseDto, err := tH.service.UpdateTaskDefinitionById(r.Context(), uId, tDId, tDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

func (tH *TaskHandler) deleteTaskDefinition(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDId, err := uuid.Parse(chi.URLParam(r, "tDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDResponseDto, err := tH.service.SoftDeleteTaskDefinition(r.Context(), uId, tDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDResponseDto)
}

// TaskDependency

func (tH *TaskHandler) createTaskDependency(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	var tDpRequestDto task.TaskDependencyRequestDto
	err := render.DecodeJSON(r.Body, &tDpRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := tH.service.CreateTaskDependency(r.Context(), uId, tDpRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (tH *TaskHandler) getTaskDependencyByWorkflowDefinition(
	w http.ResponseWriter,
	r *http.Request,
) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wDId, err := uuid.Parse(chi.URLParam(r, "wDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := tH.service.GetTaskDependencyByWorkflowDefinition(r.Context(), uId, wDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (tH *TaskHandler) getTaskDependencyByUserAndId(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDpId, err := uuid.Parse(chi.URLParam(r, "tDpId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := tH.service.GetTaskDependencyByUserAndId(r.Context(), uId, tDpId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (tH *TaskHandler) updateTaskDependencyById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDpId, err := uuid.Parse(chi.URLParam(r, "tDpId"))
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

	tDpResponseDto, err := tH.service.UpdateTaskDependencyById(
		r.Context(),
		uId,
		tDpId,
		tDpRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

func (tH *TaskHandler) deleteTaskDependency(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tDpId, err := uuid.Parse(chi.URLParam(r, "tDpId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDpResponseDto, err := tH.service.SoftDeleteTaskDependency(r.Context(), uId, tDpId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tDpResponseDto)
}

// TaskRun

func (tH *TaskHandler) createTaskRun(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wRId, err := uuid.Parse(chi.URLParam(r, "wRId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tDId, err := uuid.Parse(chi.URLParam(r, "tDId"))
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

	tRResponseDto, err := tH.service.CreateTaskRun(r.Context(), uId, wRId, tDId, tRRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (tH *TaskHandler) getTaskRunByWorkflowRun(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wRId, err := uuid.Parse(chi.URLParam(r, "wRId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := tH.service.GetTaskRunByWorkflowRun(r.Context(), uId, wRId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (tH *TaskHandler) getTaskRunByUserAndId(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tRId, err := uuid.Parse(chi.URLParam(r, "tRId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := tH.service.GetTaskRunByUserAndId(r.Context(), uId, tRId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (tH *TaskHandler) updateTaskRunById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tRId, err := uuid.Parse(chi.URLParam(r, "tRId"))
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

	tRResponseDto, err := tH.service.UpdateTaskRunByUserAndId(r.Context(), uId, tRId, tRRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

func (tH *TaskHandler) deleteTaskRun(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tRId, err := uuid.Parse(chi.URLParam(r, "tRId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tRResponseDto, err := tH.service.SoftDeleteTaskRun(r.Context(), uId, tRId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tRResponseDto)
}

// TaskAttempt

func (tH *TaskHandler) getTaskAttemptByTaskRun(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tRId, err := uuid.Parse(chi.URLParam(r, "tRId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := tH.service.GetTaskAttemptByTaskRun(r.Context(), uId, tRId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}

func (tH *TaskHandler) getTaskAttemptByUserAndId(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tAId, err := uuid.Parse(chi.URLParam(r, "tAId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := tH.service.GetTaskAttemptByUserAndId(r.Context(), uId, tAId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}

func (tH *TaskHandler) deleteTaskAttempt(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	tAId, err := uuid.Parse(chi.URLParam(r, "tAId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tAResponseDto, err := tH.service.SoftDeleteTaskAttempt(r.Context(), uId, tAId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, tAResponseDto)
}
