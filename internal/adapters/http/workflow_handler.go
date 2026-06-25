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

		r.Route("/{wId}", func(r chi.Router) {
			r.Get("", wH.getWorkflowByUserAndId)
			r.Put("", wH.updateWorkflowById)
			r.Delete("", wH.deleteWorkflowById)

			r.Route("/workflow_definition", func(r chi.Router) {
				r.Post("", wH.createWorkflowDefinition)
				r.Get("", wH.getWorkflowDefinitionByUserAndWorkflow)
			})
		})
	})

	r.Route("/workflow_definition/{wDId}", func(r chi.Router) {
		r.Get("", wH.getWorkflowDefinitionByUserAndId)
		r.Put("", wH.updateWorkflowDefinitionById)
		r.Delete("", wH.deleteWorkflowDefinitionById)
	})

	return r
}

func (wH *WorkflowHandler) createWorkflow(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	var wRequestDto workflow.WorkflowRequestDto
	err := render.DecodeJSON(r.Body, &wRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wResponseDto, err := wH.service.CreateWorkflow(r.Context(), uId, wRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wResponseDto)
}

func (wH *WorkflowHandler) getWorkflowByUser(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wResponseDto, err := wH.service.GetWorkflowByUser(r.Context(), uId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wResponseDto)
}

func (wH *WorkflowHandler) getWorkflowByUserAndId(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wId, err := uuid.Parse(chi.URLParam(r, "wId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	wResponseDto, err := wH.service.GetWorkflowByUserAndId(r.Context(), uId, wId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wResponseDto)
}

func (wH *WorkflowHandler) updateWorkflowById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wId, err := uuid.Parse(chi.URLParam(r, "wId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	var wRequestDto workflow.WorkflowRequestDto
	err = render.DecodeJSON(r.Body, &wRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wResponseDto, err := wH.service.UpdateWorkflowById(r.Context(), uId, wId, wRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wResponseDto)
}

func (wH *WorkflowHandler) deleteWorkflowById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wId, err := uuid.Parse(chi.URLParam(r, "wId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	wResponseDto, err := wH.service.SoftDeleteWorkflow(r.Context(), uId, wId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wResponseDto)
}

// WorkflowDerfinition

func (wH *WorkflowHandler) createWorkflowDefinition(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wId, err := uuid.Parse(chi.URLParam(r, "wId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	var wDRequestDto workflow.WorkflowDefinitionRequestDto
	err = render.DecodeJSON(r.Body, &wDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.CreateWorkflowDefinition(r.Context(), uId, wId)
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
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wId, err := uuid.Parse(chi.URLParam(r, "wId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	wDResponseDto, err := wH.service.GetWorkflowDefinitionByUserAndWorkflow(r.Context(), uId, wId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) getWorkflowDefinitionByUserAndId(
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
	}

	wDResponseDto, err := wH.service.GetWorkflowDefinitionByUserAndId(r.Context(), uId, wDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) updateWorkflowDefinitionById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wDId, err := uuid.Parse(chi.URLParam(r, "wDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	var wDRequestDto workflow.WorkflowDefinitionRequestDto
	err = render.DecodeJSON(r.Body, &wDRequestDto)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	wDResponseDto, err := wH.service.UpdateWorkflowDefinitionById(
		r.Context(),
		uId,
		wDId,
		wDRequestDto,
	)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}

func (wH *WorkflowHandler) deleteWorkflowDefinitionById(w http.ResponseWriter, r *http.Request) {
	uId, valid := userIDFromToken(w, r)
	if !valid {
		return
	}

	wDId, err := uuid.Parse(chi.URLParam(r, "wDId"))
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
	}

	wDResponseDto, err := wH.service.SoftDeleteWorkflowDefinition(r.Context(), uId, wDId)
	if err != nil {
		renderError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	renderJSON(w, r, http.StatusOK, wDResponseDto)
}
