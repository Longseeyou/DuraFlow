package http

import (
	"net/http"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type WorkflowHandler struct {
	service *workflow.WorkflowService
}

func NewWorkflowHandler(service workflow.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{service: &service}
}

func (wH *WorkflowHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}

func (wH *WorkflowHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/workflow", func(r chi.Router) {
		r.Post("", wH.createWorkflow)
		r.Get("", wH.getWorkflowByUser)

		r.Route("/{workflowID}", func(r chi.Router) {
			r.Get("", wH.getWorkflowByUserAndID)
			r.Put("", wH.updateWorkflowByUserAndID)
			r.Delete("", wH.deleteWorkflow)

			r.Route("/workflow_definition", func(r chi.Router) {
				r.Post("", wH.createWorkflowDefinition)
				r.Get("", wH.getWorkflowDefinitionByUserAndWorkflow)
			})
		})
	})

	r.Route("/workflow_definition/{workflowDefinitionID}", func(r chi.Router) {
		r.Get("", wH.getWorkflowDefinitionByUserAndID)
		r.Put("", wH.updateWorkflowDefinitionByUserAndID)
		r.Delete("", wH.deleteWorkflowDefinition)

		r.Route("/workflow_run", func(r chi.Router) {
			r.Get("", wH.getWorkflowRunByWorkflowDefinition)
		})
	})

	r.Route("/workflow_run/{workflowRunID}", func(r chi.Router) {
		r.Get("", wH.getWorkflowRunByUserAndID)
		r.Delete("", wH.deleteWorkflowRun)
	})

	return r
}

func (wH *WorkflowHandler) createWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	var workflowRunequestDto workflow.WorkflowRequestDto
	err := render.DecodeJSON(r.Body, &workflowRunequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	workflowRunesponseDto, err := wH.service.CreateWorkflow(
		r.Context(),
		userID,
		workflowRunequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunesponseDto)
}

func (wH *WorkflowHandler) getWorkflowByUser(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowRunesponseDto, err := wH.service.GetWorkflowByUser(r.Context(), userID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunesponseDto)
}

func (wH *WorkflowHandler) getWorkflowByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	workflowRunesponseDto, err := wH.service.GetWorkflowByUserAndID(r.Context(), userID, workflowID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunesponseDto)
}

func (wH *WorkflowHandler) updateWorkflowByUserAndID(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var workflowRunequestDto workflow.WorkflowRequestDto
	err = render.DecodeJSON(r.Body, &workflowRunequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	workflowRunesponseDto, err := wH.service.UpdateWorkflowByUserAndID(
		r.Context(),
		userID,
		workflowID,
		workflowRunequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunesponseDto)
}

func (wH *WorkflowHandler) deleteWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	workflowRunesponseDto, err := wH.service.SoftDeleteWorkflow(r.Context(), userID, workflowID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunesponseDto)
}

// WorkflowDerfinition

func (wH *WorkflowHandler) createWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.CreateWorkflowDefinition(r.Context(), userID, workflowID)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) getWorkflowDefinitionByUserAndWorkflow(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.GetWorkflowDefinitionByUserAndWorkflow(
		r.Context(),
		userID,
		workflowID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) getWorkflowDefinitionByUserAndID(
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

	wDResponseDto, err := wH.service.GetWorkflowDefinitionByUserAndID(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) updateWorkflowDefinitionByUserAndID(
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

	var wDRequestDto workflow.WorkflowDefinitionRequestDto
	err = render.DecodeJSON(r.Body, &wDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.UpdateWorkflowDefinitionByUserAndID(
		r.Context(),
		userID,
		workflowDefinitionID,
		wDRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) deleteWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowDefinitionID, err := uuid.Parse(chi.URLParam(r, "workflowDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.SoftDeleteWorkflowDefinition(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

// WorkflowRun

func (wH *WorkflowHandler) getWorkflowRunByWorkflowDefinition(
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

	workflowRunResponseDto, err := wH.service.GetWorkflowRunByWorkflowDefinition(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunResponseDto)
}

func (wH *WorkflowHandler) getWorkflowRunByUserAndID(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowRunID, err := uuid.Parse(chi.URLParam(r, "workflowRunID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	workflowRunResponseDto, err := wH.service.GetWorkflowRunByUserAndID(
		r.Context(),
		userID,
		workflowRunID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, workflowRunResponseDto)
}

func (wH *WorkflowHandler) deleteWorkflowRun(w http.ResponseWriter, r *http.Request) {
	userID, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	workflowDefinitionID, err := uuid.Parse(chi.URLParam(r, "workflowDefinitionID"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.SoftDeleteWorkflowRun(
		r.Context(),
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}
