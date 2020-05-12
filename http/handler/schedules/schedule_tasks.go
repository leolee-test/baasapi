package schedules

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
)

type taskContainer struct {
	ID         string               `json:"Id"`
	Baask8sID baasapi.Baask8sID `json:"Baask8sId"`
	Status     string               `json:"Status"`
	Created    float64              `json:"Created"`
	Labels     map[string]string    `json:"Labels"`
}

// GET request on /api/schedules/:id/tasks
func (handler *Handler) scheduleTasks(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
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

	if schedule.JobType != baasapi.ScriptExecutionJobType {
		return &httperror.HandlerError{http.StatusBadRequest, "Unable to retrieve schedule tasks", errors.New("This type of schedule do not have any associated tasks")}
	}

	tasks := make([]taskContainer, 0)

	for _, baask8sID := range schedule.ScriptExecutionJob.Baask8ss {
		baask8s, err := handler.Baask8sService.Baask8s(baask8sID)
		if err == baasapi.ErrObjectNotFound {
			continue
		} else if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find an baask8s with the specified identifier inside the database", err}
		}

		baask8sTasks, err := extractTasksFromContainerSnasphot(baask8s, schedule.ID)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find extract schedule tasks from baask8s snapshot", err}
		}

		tasks = append(tasks, baask8sTasks...)
	}

	return response.JSON(w, tasks)
}

func extractTasksFromContainerSnasphot(baask8s *baasapi.Baask8s, scheduleID baasapi.ScheduleID) ([]taskContainer, error) {
	baask8sTasks := make([]taskContainer, 0)
	if len(baask8s.Snapshots) == 0 {
		return baask8sTasks, nil
	}

	b, err := json.Marshal(baask8s.Snapshots[0].SnapshotRaw.Containers)
	if err != nil {
		return nil, err
	}

	var containers []taskContainer
	err = json.Unmarshal(b, &containers)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		if container.Labels["io.baasapi.schedule.id"] == strconv.Itoa(int(scheduleID)) {
			container.Baask8sID = baask8s.ID
			baask8sTasks = append(baask8sTasks, container)
		}
	}

	return baask8sTasks, nil
}
