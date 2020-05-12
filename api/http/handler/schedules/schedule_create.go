package schedules

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/cron"
)

type scheduleCreateFromFilePayload struct {
	Name           string
	Image          string
	CronExpression string
	Recurring      bool
	Baask8ss      []baasapi.Baask8sID
	File           []byte
	RetryCount     int
	RetryInterval  int
}

type scheduleCreateFromFileContentPayload struct {
	Name           string
	CronExpression string
	Recurring      bool
	Image          string
	Baask8ss      []baasapi.Baask8sID
	FileContent    string
	RetryCount     int
	RetryInterval  int
}

//baask8sJob := &baasapi.Baask8sJob{}
//baask8sSchedule = &baasapi.Schedule{
//	ID:             baasapi.ScheduleID(scheduleService.GetNextIdentifier()),
//	Name:           "system_baask8s",
//	CronExpression: "@every " + settings.Baask8sInterval,
//	Recurring:      true,
//	JobType:        baasapi.Baask8sJobType,
//	Baask8sJob:     baask8sJob,
//	Created:        time.Now().Unix(),
//}

type scheduleCreateFromPayload struct {
	Name           string
	CronExpression string
	Recurring      bool
	//JobType        baasapi.Baask8sJobType
	Baask8ss        []baasapi.Baask8sID
	//Baask8sJob     baasapi.Baask8sJob
}

func (payload *scheduleCreateFromFilePayload) Validate(r *http.Request) error {
	name, err := request.RetrieveMultiPartFormValue(r, "Name", false)
	if err != nil {
		return errors.New("Invalid schedule name")
	}

	if !govalidator.Matches(name, `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`) {
		return errors.New("Invalid schedule name format. Allowed characters are: [a-zA-Z0-9_.-]")
	}
	payload.Name = name

	image, err := request.RetrieveMultiPartFormValue(r, "Image", false)
	if err != nil {
		return errors.New("Invalid schedule image")
	}
	payload.Image = image

	cronExpression, err := request.RetrieveMultiPartFormValue(r, "CronExpression", false)
	if err != nil {
		return errors.New("Invalid cron expression")
	}
	payload.CronExpression = cronExpression

	var baask8ss []baasapi.Baask8sID
	err = request.RetrieveMultiPartFormJSONValue(r, "Baask8ss", &baask8ss, false)
	if err != nil {
		return errors.New("Invalid baask8ss")
	}
	payload.Baask8ss = baask8ss

	file, _, err := request.RetrieveMultiPartFormFile(r, "file")
	if err != nil {
		return baasapi.Error("Invalid script file. Ensure that the file is uploaded correctly")
	}
	payload.File = file

	retryCount, _ := request.RetrieveNumericMultiPartFormValue(r, "RetryCount", true)
	payload.RetryCount = retryCount

	retryInterval, _ := request.RetrieveNumericMultiPartFormValue(r, "RetryInterval", true)
	payload.RetryInterval = retryInterval

	return nil
}

func (payload *scheduleCreateFromFileContentPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid schedule name")
	}

	if !govalidator.Matches(payload.Name, `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`) {
		return errors.New("Invalid schedule name format. Allowed characters are: [a-zA-Z0-9_.-]")
	}

	if govalidator.IsNull(payload.Image) {
		return baasapi.Error("Invalid schedule image")
	}

	if govalidator.IsNull(payload.CronExpression) {
		return baasapi.Error("Invalid cron expression")
	}

	if payload.Baask8ss == nil || len(payload.Baask8ss) == 0 {
		return baasapi.Error("Invalid baask8ss payload")
	}

	if govalidator.IsNull(payload.FileContent) {
		return baasapi.Error("Invalid script file content")
	}

	if payload.RetryCount != 0 && payload.RetryInterval == 0 {
		return baasapi.Error("RetryInterval must be set")
	}

	return nil
}

func (payload *scheduleCreateFromPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return baasapi.Error("Invalid schedule name")
	}

	if !govalidator.Matches(payload.Name, `^[a-zA-Z0-9][a-zA-Z0-9_.-]+$`) {
		return errors.New("Invalid schedule name format. Allowed characters are: [a-zA-Z0-9_.-]")
	}

	if govalidator.IsNull(payload.CronExpression) {
		return baasapi.Error("Invalid cron expression")
	}

	if payload.Baask8ss == nil || len(payload.Baask8ss) == 0 {
		return baasapi.Error("Invalid baask8ss payload")
	}

	return nil
}

// POST /api/schedules?method=file/string/baask8s
func (handler *Handler) scheduleCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Unable to retrieve settings", err}
	}
	if !settings.EnableHostManagementFeatures {
		return &httperror.HandlerError{http.StatusServiceUnavailable, "Host management features are disabled", baasapi.ErrHostManagementFeaturesDisabled}
	}

	method, err := request.RetrieveQueryParameter(r, "method", false)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: method. Valid values are: file or string", err}
	}

	switch method {
	case "string":
		return handler.createScheduleFromFileContent(w, r)
	case "file":
		return handler.createScheduleFromFile(w, r)
	case "baask8s":
		return handler.createScheduleFromPayload(w, r)
	default:
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: method. Valid values are: file or string", errors.New(request.ErrInvalidQueryParameter)}
	}
}

func (handler *Handler) createScheduleFromPayload(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload scheduleCreateFromPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	schedule := handler.createScheduleObjectFromPayload(&payload)

	err = handler.addAndPersistScheduleBaask8s(schedule)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to schedule script job", err}
	}

	return response.JSON(w, schedule)
}

func (handler *Handler) createScheduleFromFileContent(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload scheduleCreateFromFileContentPayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	schedule := handler.createScheduleObjectFromFileContentPayload(&payload)

	err = handler.addAndPersistSchedule(schedule, []byte(payload.FileContent))
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to schedule script job", err}
	}

	return response.JSON(w, schedule)
}

func (handler *Handler) createScheduleFromFile(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	payload := &scheduleCreateFromFilePayload{}
	err := payload.Validate(r)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	schedule := handler.createScheduleObjectFromFilePayload(payload)

	err = handler.addAndPersistSchedule(schedule, payload.File)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to schedule script job", err}
	}

	return response.JSON(w, schedule)
}

func (handler *Handler) createScheduleObjectFromPayload(payload *scheduleCreateFromPayload) *baasapi.Schedule {
	scheduleIdentifier := baasapi.ScheduleID(handler.ScheduleService.GetNextIdentifier())

	job := &baasapi.Baask8sJob{
		Baask8ss:     payload.Baask8ss,
	}
	schedule := &baasapi.Schedule{
		ID:                 scheduleIdentifier,
		Name:               payload.Name,
		CronExpression:     payload.CronExpression,
		Recurring:          payload.Recurring,
		JobType:            baasapi.Baask8sJobType,
		Baask8sJob:         job,
		Created:            time.Now().Unix(),
	}

	return schedule
}

func (handler *Handler) createScheduleObjectFromFilePayload(payload *scheduleCreateFromFilePayload) *baasapi.Schedule {
	scheduleIdentifier := baasapi.ScheduleID(handler.ScheduleService.GetNextIdentifier())

	job := &baasapi.ScriptExecutionJob{
		Baask8ss:     payload.Baask8ss,
		Image:         payload.Image,
		RetryCount:    payload.RetryCount,
		RetryInterval: payload.RetryInterval,
	}

	schedule := &baasapi.Schedule{
		ID:                 scheduleIdentifier,
		Name:               payload.Name,
		CronExpression:     payload.CronExpression,
		Recurring:          payload.Recurring,
		JobType:            baasapi.ScriptExecutionJobType,
		ScriptExecutionJob: job,
		Created:            time.Now().Unix(),
	}

	return schedule
}

func (handler *Handler) createScheduleObjectFromFileContentPayload(payload *scheduleCreateFromFileContentPayload) *baasapi.Schedule {
	scheduleIdentifier := baasapi.ScheduleID(handler.ScheduleService.GetNextIdentifier())

	job := &baasapi.ScriptExecutionJob{
		Baask8ss:     payload.Baask8ss,
		Image:         payload.Image,
		RetryCount:    payload.RetryCount,
		RetryInterval: payload.RetryInterval,
	}

	schedule := &baasapi.Schedule{
		ID:                 scheduleIdentifier,
		Name:               payload.Name,
		CronExpression:     payload.CronExpression,
		Recurring:          payload.Recurring,
		JobType:            baasapi.ScriptExecutionJobType,
		ScriptExecutionJob: job,
		Created:            time.Now().Unix(),
	}

	return schedule
}

func (handler *Handler) addAndPersistSchedule(schedule *baasapi.Schedule, file []byte) error {
	scriptPath, err := handler.FileService.StoreScheduledJobFileFromBytes(strconv.Itoa(int(schedule.ID)), file)
	if err != nil {
		return err
	}

	schedule.ScriptExecutionJob.ScriptPath = scriptPath

	//jobContext := cron.NewScriptExecutionJobContext(handler.JobService, handler.Baask8sService, handler.FileService)
	//jobRunner := cron.NewScriptExecutionJobRunner(schedule, jobContext)

	//err = handler.JobScheduler.ScheduleJob(jobRunner)
	//if err != nil {
	//	return err
	//}

	return handler.ScheduleService.CreateSchedule(schedule)
}

func (handler *Handler) addAndPersistScheduleBaask8s(schedule *baasapi.Schedule) error {

	jobContext := cron.NewBaask8sJobContext(handler.Baask8sService, handler.CAFilesManager)
	jobRunner := cron.NewBaask8sJobRunner(schedule, jobContext)

	err := handler.JobScheduler.ScheduleJob(jobRunner)
	if err != nil {
		return err
	}

	return handler.ScheduleService.CreateSchedule(schedule)
}
