package status

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
)

// GET request on /api/status
func (handler *Handler) statusInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	return response.JSON(w, handler.Status)
}
