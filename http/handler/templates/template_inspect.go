package templates

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// GET request on /api/templates/:id
func (handler *Handler) templateInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	templateID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid template identifier route variable", err}
	}

	template, err := handler.TemplateService.Template(baasapi.TemplateID(templateID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a template with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a template with the specified identifier inside the database", err}
	}

	return response.JSON(w, template)
}
