package settings

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/libhttp/response"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/filesystem"
)

type settingsUpdatePayload struct {
	LogoURL                            *string
	BlackListedLabels                  []baasapi.Pair
	AuthenticationMethod               *int
	LDAPSettings                       *baasapi.LDAPSettings
	OAuthSettings                      *baasapi.OAuthSettings
	AllowBindMountsForRegularUsers     *bool
	AllowPrivilegedModeForRegularUsers *bool
	EnableHostManagementFeatures       *bool
	SnapshotInterval                   *string
	Baask8sInterval                    *string
	TemplatesURL                       *string
}

func (payload *settingsUpdatePayload) Validate(r *http.Request) error {
	if *payload.AuthenticationMethod != 1 && *payload.AuthenticationMethod != 2 && *payload.AuthenticationMethod != 3 {
		return baasapi.Error("Invalid authentication method value. Value must be one of: 1 (internal), 2 (LDAP/AD) or 3 (OAuth)")
	}
	if payload.LogoURL != nil && *payload.LogoURL != "" && !govalidator.IsURL(*payload.LogoURL) {
		return baasapi.Error("Invalid logo URL. Must correspond to a valid URL format")
	}
	if payload.TemplatesURL != nil && *payload.TemplatesURL != "" && !govalidator.IsURL(*payload.TemplatesURL) {
		return baasapi.Error("Invalid external templates URL. Must correspond to a valid URL format")
	}
	return nil
}

// PUT request on /api/settings
func (handler *Handler) settingsUpdate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload settingsUpdatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	settings, err := handler.SettingsService.Settings()
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to retrieve the settings from the database", err}
	}

	if payload.AuthenticationMethod != nil {
		settings.AuthenticationMethod = baasapi.AuthenticationMethod(*payload.AuthenticationMethod)
	}

	if payload.LogoURL != nil {
		settings.LogoURL = *payload.LogoURL
	}

	if payload.TemplatesURL != nil {
		settings.TemplatesURL = *payload.TemplatesURL
	}

	if payload.BlackListedLabels != nil {
		settings.BlackListedLabels = payload.BlackListedLabels
	}

	if payload.LDAPSettings != nil {
		ldapPassword := settings.LDAPSettings.Password
		if payload.LDAPSettings.Password != "" {
			ldapPassword = payload.LDAPSettings.Password
		}
		settings.LDAPSettings = *payload.LDAPSettings
		settings.LDAPSettings.Password = ldapPassword
	}

	if payload.OAuthSettings != nil {
		clientSecret := payload.OAuthSettings.ClientSecret
		if clientSecret == "" {
			clientSecret = settings.OAuthSettings.ClientSecret
		}
		settings.OAuthSettings = *payload.OAuthSettings
		settings.OAuthSettings.ClientSecret = clientSecret
	}

	if payload.AllowBindMountsForRegularUsers != nil {
		settings.AllowBindMountsForRegularUsers = *payload.AllowBindMountsForRegularUsers
	}

	if payload.AllowPrivilegedModeForRegularUsers != nil {
		settings.AllowPrivilegedModeForRegularUsers = *payload.AllowPrivilegedModeForRegularUsers
	}

	if payload.EnableHostManagementFeatures != nil {
		settings.EnableHostManagementFeatures = *payload.EnableHostManagementFeatures
	}

	if payload.SnapshotInterval != nil && *payload.SnapshotInterval != settings.SnapshotInterval {
		err := handler.updateSnapshotInterval(settings, *payload.SnapshotInterval)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update snapshot interval", err}
		}
	}

	if payload.Baask8sInterval != nil && *payload.Baask8sInterval != settings.Baask8sInterval {
		err := handler.updateBaask8sInterval(settings, *payload.Baask8sInterval)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to update baask8s interval", err}
		}
	}

	tlsError := handler.updateTLS(settings)
	if tlsError != nil {
		return tlsError
	}

	err = handler.SettingsService.UpdateSettings(settings)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to persist settings changes inside the database", err}
	}

	return response.JSON(w, settings)
}

func (handler *Handler) updateSnapshotInterval(settings *baasapi.Settings, snapshotInterval string) error {
	settings.SnapshotInterval = snapshotInterval

	schedules, err := handler.ScheduleService.SchedulesByJobType(baasapi.SnapshotJobType)
	if err != nil {
		return err
	}

	if len(schedules) != 0 {
		snapshotSchedule := schedules[0]
		snapshotSchedule.CronExpression = "@every " + snapshotInterval

		err := handler.JobScheduler.UpdateSystemJobSchedule(baasapi.SnapshotJobType, snapshotSchedule.CronExpression)
		if err != nil {
			return err
		}

		err = handler.ScheduleService.UpdateSchedule(snapshotSchedule.ID, &snapshotSchedule)
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *Handler) updateBaask8sInterval(settings *baasapi.Settings, baask8sInterval string) error {
	settings.Baask8sInterval = baask8sInterval

	schedules, err := handler.ScheduleService.SchedulesByJobType(baasapi.Baask8sJobType)
	if err != nil {
		return err
	}

	if len(schedules) != 0 {
		baask8sSchedule := schedules[0]
		baask8sSchedule.CronExpression = "@every " + baask8sInterval

		err := handler.JobScheduler.UpdateSystemJobSchedule(baasapi.Baask8sJobType, baask8sSchedule.CronExpression)
		if err != nil {
			return err
		}

		err = handler.ScheduleService.UpdateSchedule(baask8sSchedule.ID, &baask8sSchedule)
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *Handler) updateTLS(settings *baasapi.Settings) *httperror.HandlerError {
	if (settings.LDAPSettings.TLSConfig.TLS || settings.LDAPSettings.StartTLS) && !settings.LDAPSettings.TLSConfig.TLSSkipVerify {
		caCertPath, _ := handler.FileService.GetPathForTLSFile(filesystem.LDAPStorePath, baasapi.TLSFileCA)
		settings.LDAPSettings.TLSConfig.TLSCACertPath = caCertPath
	} else {
		settings.LDAPSettings.TLSConfig.TLSCACertPath = ""
		err := handler.FileService.DeleteTLSFiles(filesystem.LDAPStorePath)
		if err != nil {
			return &httperror.HandlerError{http.StatusInternalServerError, "Unable to remove TLS files from disk", err}
		}
	}
	return nil
}
