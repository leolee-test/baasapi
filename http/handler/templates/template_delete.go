package templates

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/templates/:id
func (handler *Handler) templateDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid template identifier route variable", err}
	}

	err = handler.TemplateService.DeleteTemplate(baasapi.TemplateID(id))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove the template from the database", err}
	}

	return response.Empty(w)
}
