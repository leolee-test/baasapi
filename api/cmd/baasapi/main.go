package main

import (
	"encoding/json"
	"os"
	//"strings"
	"time"

	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/bolt"
	"github.com/baasapi/baasapi/api/cli"
	"github.com/baasapi/baasapi/api/cron"
	"github.com/baasapi/baasapi/api/crypto"
	//"github.com/baasapi/baasapi/api/docker"
	"github.com/baasapi/baasapi/api/exec"
	"github.com/baasapi/baasapi/api/filesystem"
	"github.com/baasapi/baasapi/api/git"
	"github.com/baasapi/baasapi/api/http"
	//"github.com/baasapi/baasapi/api/http/client"
	"github.com/baasapi/baasapi/api/jwt"
	"github.com/baasapi/baasapi/api/ldap"
	//"github.com/baasapi/baasapi/api/libcompose"

	"log"
)

func initCLI() *baasapi.CLIFlags {
	var cli baasapi.CLIService = &cli.Service{}
	flags, err := cli.ParseFlags(baasapi.APIVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = cli.ValidateFlags(flags)
	if err != nil {
		log.Fatal(err)
	}
	return flags
}

func initFileService(dataStorePath string) baasapi.FileService {
	fileService, err := filesystem.NewService(dataStorePath, "")
	if err != nil {
		log.Fatal(err)
	}
	return fileService
}

func initStore(dataStorePath string, fileService baasapi.FileService) *bolt.Store {
	store, err := bolt.NewStore(dataStorePath, fileService)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Open()
	if err != nil {
		log.Fatal(err)
	}

	err = store.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = store.MigrateData()
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func initCAFilesManager(assetsPath string, dataStorePath string, fileService baasapi.FileService) (baasapi.CAFilesManager, error) {
	return exec.NewCAFilesManager(assetsPath, dataStorePath, fileService)
}

func initJWTService(authenticationEnabled bool) baasapi.JWTService {
	if authenticationEnabled {
		jwtService, err := jwt.NewService()
		if err != nil {
			log.Fatal(err)
		}
		return jwtService
	}
	return nil
}

func initDigitalSignatureService() baasapi.DigitalSignatureService {
	return crypto.NewECDSAService(os.Getenv("AGENT_SECRET"))
}

func initCryptoService() baasapi.CryptoService {
	return &crypto.Service{}
}

func initLDAPService() baasapi.LDAPService {
	return &ldap.Service{}
}

func initGitService() baasapi.GitService {
	return &git.Service{}
}

func initJobScheduler() baasapi.JobScheduler {
	return cron.NewJobScheduler()
}

//func loadSnapshotSystemSchedule(jobScheduler baasapi.JobScheduler, snapshotter baasapi.Snapshotter, scheduleService baasapi.ScheduleService, baask8sService baasapi.Baask8sService, settingsService baasapi.SettingsService) error {
//	settings, err := settingsService.Settings()
//	if err != nil {
//		return err
//	}

//	schedules, err := scheduleService.SchedulesByJobType(baasapi.SnapshotJobType)
//	if err != nil {
//		return err
//	}

//	var snapshotSchedule *baasapi.Schedule
//	if len(schedules) == 0 {
//		snapshotJob := &baasapi.SnapshotJob{}
//		snapshotSchedule = &baasapi.Schedule{
//			ID:             baasapi.ScheduleID(scheduleService.GetNextIdentifier()),
//			Name:           "system_snapshot",
//			CronExpression: "@every " + settings.SnapshotInterval,
//			Recurring:      true,
//			JobType:        baasapi.SnapshotJobType,
//			SnapshotJob:    snapshotJob,
//			Created:        time.Now().Unix(),
//		}
//	} else {
//		snapshotSchedule = &schedules[0]
//	}
//
//	snapshotJobContext := cron.NewSnapshotJobContext(baask8sService, snapshotter)
//	snapshotJobRunner := cron.NewSnapshotJobRunner(snapshotSchedule, snapshotJobContext)
//
//	err = jobScheduler.ScheduleJob(snapshotJobRunner)
//	if err != nil {
//		return err
//	}
//
//	if len(schedules) == 0 {
//		return scheduleService.CreateSchedule(snapshotSchedule)
//	}
//	return nil
//}

func loadBaask8sSystemSchedule(jobScheduler baasapi.JobScheduler, scheduleService baasapi.ScheduleService, baask8sService baasapi.Baask8sService, settingsService baasapi.SettingsService, cafilesmanager baasapi.CAFilesManager) error {
	
	//schedules, err := scheduleService.Schedules()
	//if err != nil {
	//	return err
	//}

	//for _, schedule := range schedules {
	//	if schedule.JobType == baasapi.Baask8sJobType {
	//		jobContext := cron.NewBaask8sJobContext(baask8sService, cafilesmanager)
	//		jobRunner := cron.NewBaask8sJobRunner(&schedule, jobContext)
//
//			err = jobScheduler.ScheduleJob(jobRunner)
//			if err != nil {
//				return err
//			}
//		}
//	}

//	return nil
	
	settings, err := settingsService.Settings()
	if err != nil {
		return err
	}

	schedules, err := scheduleService.SchedulesByJobType(baasapi.Baask8sJobType)
	if err != nil {
		return err
	}

	var baask8sSchedule *baasapi.Schedule
	if len(schedules) == 0 {
		baask8sJob := &baasapi.Baask8sJob{}
		baask8sSchedule = &baasapi.Schedule{
			ID:             baasapi.ScheduleID(scheduleService.GetNextIdentifier()),
			Name:           "system_baask8s",
			CronExpression: "@every " + settings.Baask8sInterval,
			Recurring:      true,
			JobType:        baasapi.Baask8sJobType,
			Baask8sJob:     baask8sJob,
			Created:        time.Now().Unix(),
		}
	} else {
		baask8sSchedule = &schedules[0]
	}

	baask8sJobContext := cron.NewBaask8sJobContext(baask8sService, cafilesmanager)
    baask8sJobRunner := cron.NewBaask8sJobRunner(baask8sSchedule, baask8sJobContext)

	err = jobScheduler.ScheduleJob(baask8sJobRunner)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		return scheduleService.CreateSchedule(baask8sSchedule)
	}
	return nil
	
}

//func loadSchedulesFromDatabase(jobScheduler baasapi.JobScheduler, jobService baasapi.JobService, scheduleService baasapi.ScheduleService, baask8sService baasapi.Baask8sService, fileService baasapi.FileService) error {
//	schedules, err := scheduleService.Schedules()
//	if err != nil {
//		return err
//	}

//	for _, schedule := range schedules {

		//if schedule.JobType == baasapi.ScriptExecutionJobType {
		//	jobContext := cron.NewScriptExecutionJobContext(jobService, baask8sService, fileService)
		//	jobRunner := cron.NewScriptExecutionJobRunner(&schedule, jobContext)
//
//			err = jobScheduler.ScheduleJob(jobRunner)
//			if err != nil {
//				return err
//			}
//		}
//	}

//	return nil
//}

func initStatus(baask8sManagement, snapshot bool, flags *baasapi.CLIFlags) *baasapi.Status {
	return &baasapi.Status{
		Analytics:          !*flags.NoAnalytics,
		Authentication:     !*flags.NoAuth,
		Baask8sManagement: baask8sManagement,
		Snapshot:           snapshot,
		Version:            baasapi.APIVersion,
	}
}

func initSettings(settingsService baasapi.SettingsService, flags *baasapi.CLIFlags) error {
	_, err := settingsService.Settings()
	if err == baasapi.ErrObjectNotFound {
		settings := &baasapi.Settings{
			LogoURL:              *flags.Logo,
			AuthenticationMethod: baasapi.AuthenticationInternal,
			LDAPSettings: baasapi.LDAPSettings{
				AutoCreateUsers: true,
				TLSConfig:       baasapi.TLSConfiguration{},
				SearchSettings: []baasapi.LDAPSearchSettings{
					baasapi.LDAPSearchSettings{},
				},
				GroupSearchSettings: []baasapi.LDAPGroupSearchSettings{
					baasapi.LDAPGroupSearchSettings{},
				},
			},
			OAuthSettings:                      baasapi.OAuthSettings{},
			AllowBindMountsForRegularUsers:     true,
			AllowPrivilegedModeForRegularUsers: true,
			EnableHostManagementFeatures:       true,
			SnapshotInterval:                   *flags.SnapshotInterval,
			Baask8sInterval:                    *flags.Baask8sInterval,
		}

		if *flags.Templates != "" {
			settings.TemplatesURL = *flags.Templates
		}

		if *flags.Labels != nil {
			settings.BlackListedLabels = *flags.Labels
		} else {
			settings.BlackListedLabels = make([]baasapi.Pair, 0)
		}

		return settingsService.UpdateSettings(settings)
	} else if err != nil {
		return err
	}

	return nil
}

func initTemplates(templateService baasapi.TemplateService, fileService baasapi.FileService, templateURL, templateFile string) error {
	if templateURL != "" {
		log.Printf("BaaSapi started with the --templates flag. Using external templates, template management will be disabled.")
		return nil
	}

	existingTemplates, err := templateService.Templates()
	if err != nil {
		return err
	}

	if len(existingTemplates) != 0 {
		log.Printf("Templates already registered inside the database. Skipping template import.")
		return nil
	}

	templatesJSON, err := fileService.GetFileContent(templateFile)
	if err != nil {
		log.Println("Unable to retrieve template definitions via filesystem")
		return err
	}

	var templates []baasapi.Template
	err = json.Unmarshal(templatesJSON, &templates)
	if err != nil {
		log.Println("Unable to parse templates file. Please review your template definition file.")
		return err
	}

	for _, template := range templates {
		err := templateService.CreateTemplate(&template)
		if err != nil {
			return err
		}
	}

	return nil
}


func loadAndParseKeyPair(fileService baasapi.FileService, signatureService baasapi.DigitalSignatureService) error {
	private, public, err := fileService.LoadKeyPair()
	if err != nil {
		return err
	}
	return signatureService.ParseKeyPair(private, public)
}

func generateAndStoreKeyPair(fileService baasapi.FileService, signatureService baasapi.DigitalSignatureService) error {
	private, public, err := signatureService.GenerateKeyPair()
	if err != nil {
		return err
	}
	privateHeader, publicHeader := signatureService.PEMHeaders()
	return fileService.StoreKeyPair(private, public, privateHeader, publicHeader)
}

//func initClientFactory(signatureService baasapi.DigitalSignatureService) *docker.ClientFactory {
//	return docker.NewClientFactory(signatureService)
//}

func initKeyPair(fileService baasapi.FileService, signatureService baasapi.DigitalSignatureService) error {
	existingKeyPair, err := fileService.KeyPairFilesExist()
	if err != nil {
		log.Fatal(err)
	}

	if existingKeyPair {
		return loadAndParseKeyPair(fileService, signatureService)
	}
	return generateAndStoreKeyPair(fileService, signatureService)
}

func initExtensionManager(fileService baasapi.FileService, extensionService baasapi.ExtensionService) (baasapi.ExtensionManager, error) {
	extensionManager := exec.NewExtensionManager(fileService, extensionService)

	extensions, err := extensionService.Extensions()
	if err != nil {
		return nil, err
	}

	for _, extension := range extensions {
		err := extensionManager.EnableExtension(&extension, extension.License.LicenseKey)
		if err != nil {
			log.Printf("Unable to enable extension: %s [extension: %s]", err.Error(), extension.Name)
			extension.Enabled = false
			extension.License.Valid = false
			extensionService.Persist(&extension)
		}
	}

	return extensionManager, nil
}

func terminateIfNoAdminCreated(userService baasapi.UserService) {
	timer1 := time.NewTimer(5 * time.Minute)
	<-timer1.C

	users, err := userService.UsersByRole(baasapi.AdministratorRole)
	if err != nil {
		log.Fatal(err)
	}

	if len(users) == 0 {
		log.Fatal("No administrator account was created after 5 min. Shutting down the BaaSapi instance for security reasons.")
		return
	}
}

func main() {
	flags := initCLI()

	fileService := initFileService(*flags.Data)

	store := initStore(*flags.Data, fileService)
	defer store.Close()

	jwtService := initJWTService(!*flags.NoAuth)

	ldapService := initLDAPService()

	gitService := initGitService()

	cryptoService := initCryptoService()

	digitalSignatureService := initDigitalSignatureService()

	err := initKeyPair(fileService, digitalSignatureService)
	if err != nil {
		log.Fatal(err)
	}

	extensionManager, err := initExtensionManager(fileService, store.ExtensionService)
	if err != nil {
		log.Fatal(err)
	}

	//clientFactory := initClientFactory(digitalSignatureService)

	//jobService := initJobService(clientFactory)

	//snapshotter := initSnapshotter(clientFactory)

	baask8sManagement := true
	//if *flags.ExternalBaask8ss != "" {
	//	baask8sManagement = false
	//}

	caFilesManager, err := initCAFilesManager(*flags.Assets, *flags.Data, fileService)
	if err != nil {
		log.Fatal(err)
	}

	err = initTemplates(store.TemplateService, fileService, *flags.Templates, *flags.TemplateFile)
	if err != nil {
		log.Fatal(err)
	}

	err = initSettings(store.SettingsService, flags)
	if err != nil {
		log.Fatal(err)
	}

	jobScheduler := initJobScheduler()

	//err = loadSchedulesFromDatabase(jobScheduler, jobService, store.ScheduleService, store.Baask8sService, fileService)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//if *flags.Snapshot {
	//	err = loadSnapshotSystemSchedule(jobScheduler, snapshotter, store.ScheduleService, store.Baask8sService, store.SettingsService)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}

	if *flags.Baask8s {
		err = loadBaask8sSystemSchedule(jobScheduler, store.ScheduleService, store.Baask8sService, store.SettingsService, caFilesManager)
		if err != nil {
			log.Fatal(err)
		}
	}

	jobScheduler.Start()

	applicationStatus := initStatus(baask8sManagement, *flags.Snapshot, flags)

	adminPasswordHash := ""
	if *flags.AdminPasswordFile != "" {
		content, err := fileService.GetFileContent(*flags.AdminPasswordFile)
		if err != nil {
			log.Fatal(err)
		}
		adminPasswordHash, err = cryptoService.Hash(string(content))
		if err != nil {
			log.Fatal(err)
		}
	} else if *flags.AdminPassword != "" {
		adminPasswordHash = *flags.AdminPassword
	}

	if adminPasswordHash != "" {
		users, err := store.UserService.UsersByRole(baasapi.AdministratorRole)
		if err != nil {
			log.Fatal(err)
		}

		if len(users) == 0 {
			log.Printf("Creating admin user with password hash %s", adminPasswordHash)
			user := &baasapi.User{
				Username: "admin",
				Role:     baasapi.AdministratorRole,
				Password: adminPasswordHash,
			}
			err := store.UserService.CreateUser(user)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("Instance already has an administrator user defined. Skipping admin password related flags.")
		}
	}

	if !*flags.NoAuth {
		go terminateIfNoAdminCreated(store.UserService)
	}

	var server baasapi.Server = &http.Server{
		Status:                 applicationStatus,
		BindAddress:            *flags.Addr,
		AssetsPath:             *flags.Assets,
		AuthDisabled:           *flags.NoAuth,
		Baask8sManagement:      baask8sManagement,
		UserService:            store.UserService,
		TeamService:            store.TeamService,
		TeamMembershipService:  store.TeamMembershipService,
		Baask8sService:         store.Baask8sService,
		BaasmspService:         store.BaasmspService,
		ExtensionService:       store.ExtensionService,
		ResourceControlService: store.ResourceControlService,
		SettingsService:        store.SettingsService,
		RegistryService:        store.RegistryService,
		ScheduleService:        store.ScheduleService,
		TagService:             store.TagService,
		TemplateService:        store.TemplateService,
		CAFilesManager:         caFilesManager,
		ExtensionManager:       extensionManager,
		CryptoService:          cryptoService,
		JWTService:             jwtService,
		FileService:            fileService,
		LDAPService:            ldapService,
		GitService:             gitService,
		SignatureService:       digitalSignatureService,
		JobScheduler:           jobScheduler,
		//Snapshotter:            snapshotter,
		SSL:                    *flags.SSL,
		SSLCert:                *flags.SSLCert,
		SSLKey:                 *flags.SSLKey,
		//DockerClientFactory:    clientFactory,
		//JobService:             jobService,
	}

	log.Printf("Starting BaaSapi %s on %s", baasapi.APIVersion, *flags.Addr)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
