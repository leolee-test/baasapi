package schedules

import (
	"net/http"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

// GET request on /api/schedules
func (handler *Handler) scheduleList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Unable to retrieve settings", err}
	}
	if !settings.EnableHostManagementFeatures {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Host management features are disabled", baasapi.ErrHostManagementFeaturesDisabled}
	}

	schedules, err := handler.ScheduleService.Schedules()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve schedules from the database", err}
	}

	return response.JSON(w, schedules)
}
