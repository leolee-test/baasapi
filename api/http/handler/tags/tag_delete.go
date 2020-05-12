package tags

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// DELETE request on /api/tags/:id
func (handler *Handler) tagDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid tag identifier route variable", err}
	}

	err = handler.TagService.DeleteTag(baasapi.TagID(id))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove the tag from the database", err}
	}

	return response.Empty(w)
}
