package schedules

import (
	"net/http"

	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
)

func (handler *Handler) scheduleInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Unable to retrieve settings", err}
	}
	if !settings.EnableHostManagementFeatures {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Host management features are disabled", baasapi.ErrHostManagementFeaturesDisabled}
	}

	scheduleID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid schedule identifier route variable", err}
	}

	schedule, err := handler.ScheduleService.Schedule(baasapi.ScheduleID(scheduleID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find a schedule with the specified identifier inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find a schedule with the specified identifier inside the database", err}
	}

	return response.JSON(w, schedule)
}
