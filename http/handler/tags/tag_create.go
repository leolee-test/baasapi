package tags

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type tagCreatePayload struct {
	Name string
}

func (payload *tagCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid tag name")
	}
	return nil
}

// POST request on /api/tags
func (handler *Handler) tagCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload tagCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	tags, err := handler.TagService.Tags()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve tags from the database", err}
	}

	for _, tag := range tags {
		if tag.Name == payload.Name {
			return &httperror.HandlerError{http.StatusConflict, "This name is already associated to a tag", baasapi.ErrTagAlreadyExists}
		}
	}

	tag := &baasapi.Tag{
		Name: payload.Name,
	}

	err = handler.TagService.CreateTag(tag)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist the tag inside the database", err}
	}

	return response.JSON(w, tag)
}
