package baasapi

type (
	// Pair defines a key/value string pair
	Pair struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// CLIFlags represents the available flags on the CLI
	CLIFlags struct {
		Addr              *string
		AdminPassword     *string
		AdminPasswordFile *string
		Assets            *string
		Data              *string
		Labels            *[]Pair
		Logo              *string
		NoAuth            *bool
		NoAnalytics       *bool
		Templates         *string
		TemplateFile      *string
		TLS               *bool
		TLSSkipVerify     *bool
		TLSCacert         *string
		TLSCert           *string
		TLSKey            *string
		SSL               *bool
		SSLCert           *string
		SSLKey            *string
		SyncInterval      *string
		Snapshot          *bool
		SnapshotInterval  *string
		Baask8s           *bool
		Baask8sInterval   *string
	}

	// Status represents the application status
	Status struct {
		Authentication     bool   `json:"Authentication"`
		Baask8sManagement bool   `json:"Baask8sManagement"`
		Snapshot           bool   `json:"Snapshot"`
		Analytics          bool   `json:"Analytics"`
		Version            string `json:"Version"`
	}

	// LDAPSettings represents the settings used to connect to a LDAP server
	LDAPSettings struct {
		ReaderDN            string                    `json:"ReaderDN"`
		Password            string                    `json:"Password,omitempty"`
		URL                 string                    `json:"URL"`
		TLSConfig           TLSConfiguration          `json:"TLSConfig"`
		StartTLS            bool                      `json:"StartTLS"`
		SearchSettings      []LDAPSearchSettings      `json:"SearchSettings"`
		GroupSearchSettings []LDAPGroupSearchSettings `json:"GroupSearchSettings"`
		AutoCreateUsers     bool                      `json:"AutoCreateUsers"`
	}

	// OAuthSettings represents the settings used to authorize with an authorization server
	OAuthSettings struct {
		ClientID             string `json:"ClientID"`
		ClientSecret         string `json:"ClientSecret,omitempty"`
		AccessTokenURI       string `json:"AccessTokenURI"`
		AuthorizationURI     string `json:"AuthorizationURI"`
		ResourceURI          string `json:"ResourceURI"`
		RedirectURI          string `json:"RedirectURI"`
		UserIdentifier       string `json:"UserIdentifier"`
		Scopes               string `json:"Scopes"`
		OAuthAutoCreateUsers bool   `json:"OAuthAutoCreateUsers"`
		DefaultTeamID        TeamID `json:"DefaultTeamID"`
	}

	// TLSConfiguration represents a TLS configuration
	TLSConfiguration struct {
		TLS           bool   `json:"TLS"`
		TLSSkipVerify bool   `json:"TLSSkipVerify"`
		TLSCACertPath string `json:"TLSCACert,omitempty"`
		TLSCertPath   string `json:"TLSCert,omitempty"`
		TLSKeyPath    string `json:"TLSKey,omitempty"`
	}

	// LDAPSearchSettings represents settings used to search for users in a LDAP server
	LDAPSearchSettings struct {
		BaseDN            string `json:"BaseDN"`
		Filter            string `json:"Filter"`
		UserNameAttribute string `json:"UserNameAttribute"`
	}

	// LDAPGroupSearchSettings represents settings used to search for groups in a LDAP server
	LDAPGroupSearchSettings struct {
		GroupBaseDN    string `json:"GroupBaseDN"`
		GroupFilter    string `json:"GroupFilter"`
		GroupAttribute string `json:"GroupAttribute"`
	}

	// Settings represents the application settings
	Settings struct {
		LogoURL                            string               `json:"LogoURL"`
		BlackListedLabels                  []Pair               `json:"BlackListedLabels"`
		AuthenticationMethod               AuthenticationMethod `json:"AuthenticationMethod"`
		LDAPSettings                       LDAPSettings         `json:"LDAPSettings"`
		OAuthSettings                      OAuthSettings        `json:"OAuthSettings"`
		AllowBindMountsForRegularUsers     bool                 `json:"AllowBindMountsForRegularUsers"`
		AllowPrivilegedModeForRegularUsers bool                 `json:"AllowPrivilegedModeForRegularUsers"`
		SnapshotInterval                   string               `json:"SnapshotInterval"`
		Baask8sInterval                    string               `json:"Baask8sInterval"`
		TemplatesURL                       string               `json:"TemplatesURL"`
		EnableHostManagementFeatures       bool                 `json:"EnableHostManagementFeatures"`

		// Deprecated fields
		DisplayDonationHeader       bool
		DisplayExternalContributors bool
	}

	// User represents a user account
	User struct {
		ID       UserID   `json:"Id"`
		Username string   `json:"Username"`
		Password string   `json:"Password,omitempty"`
		Role     UserRole `json:"Role"`
	}

	// UserID represents a user identifier
	UserID int

	// UserRole represents the role of a user. It can be either an administrator
	// or a regular user
	UserRole int

	// AuthenticationMethod represents the authentication method used to authenticate a user
	AuthenticationMethod int

	// Team represents a list of user accounts
	Team struct {
		ID   TeamID `json:"Id"`
		Name string `json:"Name"`
	}

	// TeamID represents a team identifier
	TeamID int

	// TeamMembership represents a membership association between a user and a team
	TeamMembership struct {
		ID     TeamMembershipID `json:"Id"`
		UserID UserID           `json:"UserID"`
		TeamID TeamID           `json:"TeamID"`
		Role   MembershipRole   `json:"Role"`
	}

	// TeamMembershipID represents a team membership identifier
	TeamMembershipID int

	// MembershipRole represents the role of a user within a team
	MembershipRole int

	// TokenData represents the data embedded in a JWT token
	TokenData struct {
		ID       UserID
		Username string
		Role     UserRole
	}

	//INITPARM struct {
	//	Fabric   string     `json:"fabric"`
	//}

	// MSP represents a MSP
	MSP struct {
		ID       int      `json:"Id"`
		MSPName  string     `json:"MSPName"`
		ORGs     []MSPORG   `json:"ORGs"`
		Role     int     `json:"Role"`
	}

	// MSPID represents a MSP identifier
	//MSPID int

	// MSPRole represents the role of a msp 
	//MSPRole int

	// MSP-ORG represents a ORG identifier
	MSPORG struct {
		ORGName   string     `json:"ORGName"`
		Anchor    string     `json:"Anchor"`
		Peers     []string   `json:"Peers"`
	}

	// ORG-Peer represents a ORG identifier
	//ORGPeer struct {
	//	PeerName  string     `json:"PeerName"`
	//}

	// MSP represents a MSP
	CHL struct {
		ID       CHLID      `json:"Id"`
		CHLName  string     `json:"CHLName"`
		CreatedAt        string               `json:"CreatedAt"`
		ORGs     []MSPORG   `json:"ORGs"`
	}

	// CHLID represents a Channel identifier
	CHLID int

	CC struct {
		//ID              CCID      `json:"id"`
		ID              int        `json:"id"`
		CCName          string     `json:"chaincodeName"`
		CHLName         string     `json:"chlName"`
		Version         string     `json:"chaincodeVersion"`
		EndorsementPolicyparm        interface{}     `json:"endorsementPolicyparm"`
		InstallORGs     []MSPORG   `json:"installORGs"`
		InstantiateORGs []MSPORG   `json:"instantiateORGs"`
		ChaincodeType   string     `json:"chaincodeType"`
		Path            string     `json:"path"`
	}

	// CHLID represents a Channel identifier
	CCID int

	// MSP-ORG represents a ORG identifier
	//MSPORG struct {
	//	ORGName   string     `json:"ORGName"`
	//	Peers     []string   `json:"Peers"`
	//}

	// RegistryID represents a registry identifier
	RegistryID int

	// RegistryType represents a type of registry
	RegistryType int

	// Registry represents a Docker registry with all the info required
	// to connect to it
	Registry struct {
		ID                      RegistryID                       `json:"Id"`
		Type                    RegistryType                     `json:"Type"`
		Name                    string                           `json:"Name"`
		URL                     string                           `json:"URL"`
		Authentication          bool                             `json:"Authentication"`
		Username                string                           `json:"Username"`
		Password                string                           `json:"Password,omitempty"`
		AuthorizedUsers         []UserID                         `json:"AuthorizedUsers"`
		AuthorizedTeams         []TeamID                         `json:"AuthorizedTeams"`
		ManagementConfiguration *RegistryManagementConfiguration `json:"ManagementConfiguration"`
	}

	// RegistryManagementConfiguration represents a configuration that can be used to query
	// the registry API via the registry management extension.
	RegistryManagementConfiguration struct {
		Type           RegistryType     `json:"Type"`
		Authentication bool             `json:"Authentication"`
		Username       string           `json:"Username"`
		Password       string           `json:"Password"`
		TLSConfig      TLSConfiguration `json:"TLSConfig"`
	}

	// Baask8sID represents an baask8s identifier
	Baask8sID int

	// Baask8stype represents the type of an baask8s
	Baask8sType int

	// Baask8sStatus represents the status of an baask8s
	Baask8sStatus int

	// Baask8s represents a Docker baask8s with all the info required
	// to connect to it
	Baask8s struct {
		ID               Baask8sID           `json:"Id"`
		//GroupID          Baask8sGroupID      `json:"GroupId"`
		NetworkName      string              `json:"NetworkName"`
		NetworkID        string              `json:"NetworkID"`
		Owner            string              `json:"Owner"`
		Platform         string              `json:"Platform"`
		CreatedAt        string               `json:"CreatedAt"`
		MSPs             []MSP               `json:"MSPs"`
		CHLs             []CHL               `json:"CHLs"`
		Policys          []Policy            `json:"Policys"`
		CCs              []CC                `json:"CCs"`
		Applications     []BaasAPP           `json:"Applications"`
		AuthorizedUsers  []UserID            `json:"AuthorizedUsers"`
		AuthorizedTeams  []TeamID            `json:"AuthorizedTeams"`
		//Tags             []string            `json:"Tags"`
		Status           Baask8sStatus       `json:"Status"`
		Snapshots        []Snapshot          `json:"Snapshots"`
		Namespace        string              `json:"Namespace"`
	}

	BaasAPP struct {
		Name             string              `json:"name"` 
		Status           string              `json:"status"` 
	}

	Policy struct {
		Name                    string        `json:"name"` 
		Description             string        `json:"desc"` 
		EndorsementPolicyparm   interface{}   `json:"endorsementPolicyparm"`
	}

	// BaasmspID represents an baasmsp identifier
	BaasmspID int

	// Baasmsptype represents the type of an baasmsp
	BaasmspType int

	// BaasmspStatus represents the status of an baasmsp
	BaasmspStatus int

	// Baasmsp represents a baasmsp with all the info required
	// to connect to it
	Baasmsp struct {
		ID               BaasmspID           `json:"Id"`
		NetworkName      string              `json:"NetworkName"`
		NetworkID        string              `json:"NetworkID"`
		Owner            string              `json:"Owner"`
		Platform         string              `json:"Platform"`
		CreatedAt        int64               `json:"CreatedAt"`
		MSPs             []MSP               `json:"MSPs"`
		//Tags             []string            `json:"Tags"`
		Status           Baask8sStatus       `json:"Status"`
		Snapshots        []Snapshot          `json:"Snapshots"`
		Namespace        string              `json:"Namespace"`
	}



	// ScheduleID represents a schedule identifier.
	ScheduleID int

	// JobType represents a job type
	JobType int

	// ScriptExecutionJob represents a scheduled job that can execute a script via a privileged container
	ScriptExecutionJob struct {
		Baask8ss     []Baask8sID
		Image         string
		ScriptPath    string
		RetryCount    int
		RetryInterval int
	}

	// SnapshotJob represents a scheduled job that can create baask8s snapshots
	SnapshotJob struct{}

	// Baask8sJob represents a scheduled job that sync baask8s 
	Baask8sJob struct{
		Baask8ss     []Baask8sID
	}

	// Schedule represents a scheduled job.
	// It only contains a pointer to one of the JobRunner implementations
	// based on the JobType.
	// NOTE: The Recurring option is only used by ScriptExecutionJob at the moment
	Schedule struct {
		ID                 ScheduleID `json:"Id"`
		Name               string
		CronExpression     string
		Recurring          bool
		Created            int64
		JobType            JobType
		ScriptExecutionJob *ScriptExecutionJob
		SnapshotJob        *SnapshotJob
		Baask8sJob         *Baask8sJob
	}

	// WebhookID represents a webhook identifier.
	WebhookID int

	// WebhookType represents the type of resource a webhook is related to
	WebhookType int

	// Webhook represents a url webhook that can be used to update a service
	Webhook struct {
		ID          WebhookID   `json:"Id"`
		Token       string      `json:"Token"`
		ResourceID  string      `json:"ResourceId"`
		Baask8sID  Baask8sID  `json:"Baask8sId"`
		WebhookType WebhookType `json:"Type"`
	}

	// AzureCredentials represents the credentials used to connect to an Azure
	// environment.
	AzureCredentials struct {
		ApplicationID     string `json:"ApplicationID"`
		TenantID          string `json:"TenantID"`
		AuthenticationKey string `json:"AuthenticationKey"`
	}

	// Snapshot represents a snapshot of a specific baask8s at a specific time
	Snapshot struct {
		Time                  int64       `json:"Time"`
		DockerVersion         string      `json:"DockerVersion"`
		Swarm                 bool        `json:"Swarm"`
		TotalCPU              int         `json:"TotalCPU"`
		TotalMemory           int64       `json:"TotalMemory"`
		RunningContainerCount int         `json:"RunningContainerCount"`
		StoppedContainerCount int         `json:"StoppedContainerCount"`
		VolumeCount           int         `json:"VolumeCount"`
		ImageCount            int         `json:"ImageCount"`
		ServiceCount          int         `json:"ServiceCount"`
		StackCount            int         `json:"StackCount"`
		SnapshotRaw           SnapshotRaw `json:"SnapshotRaw"`
	}

	// SnapshotRaw represents all the information related to a snapshot as returned by the Docker API
	SnapshotRaw struct {
		Containers interface{} `json:"Containers"`
		Volumes    interface{} `json:"Volumes"`
		Networks   interface{} `json:"Networks"`
		Images     interface{} `json:"Images"`
		Info       interface{} `json:"Info"`
		Version    interface{} `json:"Version"`
	}

	// Baask8sGroupID represents an baask8s group identifier
	Baask8sGroupID int

	// Baask8sGroup represents a group of baask8ss
	Baask8sGroup struct {
		ID              Baask8sGroupID  `json:"Id"`
		Name            string          `json:"Name"`
		Description     string          `json:"Description"`
		AuthorizedUsers []UserID        `json:"AuthorizedUsers"`
		AuthorizedTeams []TeamID        `json:"AuthorizedTeams"`
		Tags            []string        `json:"Tags"`

		// Deprecated fields
		Labels []Pair `json:"Labels"`
	}

	// ResourceControlID represents a resource control identifier
	ResourceControlID int

	// ResourceControl represent a reference to a Docker resource with specific access controls
	ResourceControl struct {
		ID             ResourceControlID    `json:"Id"`
		ResourceID     string               `json:"ResourceId"`
		SubResourceIDs []string             `json:"SubResourceIds"`
		Type           ResourceControlType  `json:"Type"`
		UserAccesses   []UserResourceAccess `json:"UserAccesses"`
		TeamAccesses   []TeamResourceAccess `json:"TeamAccesses"`
		Public         bool                 `json:"Public"`

		// Deprecated fields
		// Deprecated in DBVersion == 2
		OwnerID     UserID              `json:"OwnerId,omitempty"`
		AccessLevel ResourceAccessLevel `json:"AccessLevel,omitempty"`

		// Deprecated in DBVersion == 14
		AdministratorsOnly bool `json:"AdministratorsOnly,omitempty"`
	}

	// ResourceControlType represents the type of resource associated to the resource control (volume, container, service...)
	ResourceControlType int

	// UserResourceAccess represents the level of control on a resource for a specific user
	UserResourceAccess struct {
		UserID      UserID              `json:"UserId"`
		AccessLevel ResourceAccessLevel `json:"AccessLevel"`
	}

	// TeamResourceAccess represents the level of control on a resource for a specific team
	TeamResourceAccess struct {
		TeamID      TeamID              `json:"TeamId"`
		AccessLevel ResourceAccessLevel `json:"AccessLevel"`
	}

	// TagID represents a tag identifier
	TagID int

	// Tag represents a tag that can be associated to a resource
	Tag struct {
		ID   TagID
		Name string `json:"Name"`
	}

	// TemplateID represents a template identifier
	TemplateID int

	// TemplateType represents the type of a template
	TemplateType int

	// Template represents an application template
	Template struct {
		// Mandatory container/stack fields
		ID                TemplateID   `json:"Id"`
		Type              TemplateType `json:"type"`
		Title             string       `json:"title"`
		Description       string       `json:"description"`
		AdministratorOnly bool         `json:"administrator_only"`

		// Mandatory container fields
		Image string `json:"image"`

		// Mandatory stack fields
		Repository TemplateRepository `json:"repository"`

		// Optional stack/container fields
		Name       string        `json:"name,omitempty"`
		Logo       string        `json:"logo,omitempty"`
		Env        []TemplateEnv `json:"env,omitempty"`
		Note       string        `json:"note,omitempty"`
		Platform   string        `json:"platform,omitempty"`
		Categories []string      `json:"categories,omitempty"`

		// Optional container fields
		Registry      string           `json:"registry,omitempty"`
		Command       string           `json:"command,omitempty"`
		Network       string           `json:"network,omitempty"`
		Volumes       []TemplateVolume `json:"volumes,omitempty"`
		Ports         []string         `json:"ports,omitempty"`
		Labels        []Pair           `json:"labels,omitempty"`
		Privileged    bool             `json:"privileged,omitempty"`
		Interactive   bool             `json:"interactive,omitempty"`
		RestartPolicy string           `json:"restart_policy,omitempty"`
		Hostname      string           `json:"hostname,omitempty"`
	}

	// TemplateEnv represents a template environment variable configuration
	TemplateEnv struct {
		Name        string              `json:"name"`
		Label       string              `json:"label,omitempty"`
		Description string              `json:"description,omitempty"`
		Default     string              `json:"default,omitempty"`
		Preset      bool                `json:"preset,omitempty"`
		Select      []TemplateEnvSelect `json:"select,omitempty"`
	}

	// TemplateVolume represents a template volume configuration
	TemplateVolume struct {
		Container string `json:"container"`
		Bind      string `json:"bind,omitempty"`
		ReadOnly  bool   `json:"readonly,omitempty"`
	}

	// TemplateRepository represents the git repository configuration for a template
	TemplateRepository struct {
		URL       string `json:"url"`
		StackFile string `json:"stackfile"`
	}

	// TemplateEnvSelect represents text/value pair that will be displayed as a choice for the
	// template user
	TemplateEnvSelect struct {
		Text    string `json:"text"`
		Value   string `json:"value"`
		Default bool   `json:"default"`
	}

	// ResourceAccessLevel represents the level of control associated to a resource
	ResourceAccessLevel int

	// TLSFileType represents a type of TLS file required to connect to a Docker baask8s.
	// It can be either a TLS CA file, a TLS certificate file or a TLS key file
	TLSFileType int

	// ExtensionID represents a extension identifier
	ExtensionID int

	// Extension represents a BaaSapi extension
	Extension struct {
		ID               ExtensionID        `json:"Id"`
		Enabled          bool               `json:"Enabled"`
		Name             string             `json:"Name,omitempty"`
		ShortDescription string             `json:"ShortDescription,omitempty"`
		Description      string             `json:"Description,omitempty"`
		DescriptionURL   string             `json:"DescriptionURL,omitempty"`
		Price            string             `json:"Price,omitempty"`
		PriceDescription string             `json:"PriceDescription,omitempty"`
		Deal             bool               `json:"Deal,omitempty"`
		Available        bool               `json:"Available,omitempty"`
		License          LicenseInformation `json:"License,omitempty"`
		Version          string             `json:"Version"`
		UpdateAvailable  bool               `json:"UpdateAvailable"`
		ShopURL          string             `json:"ShopURL,omitempty"`
		Images           []string           `json:"Images,omitempty"`
		Logo             string             `json:"Logo,omitempty"`
	}

	// LicenseInformation represents information about an extension license
	LicenseInformation struct {
		LicenseKey string `json:"LicenseKey,omitempty"`
		Company    string `json:"Company,omitempty"`
		Expiration string `json:"Expiration,omitempty"`
		Valid      bool   `json:"Valid,omitempty"`
	}

	// CLIService represents a service for managing CLI
	CLIService interface {
		ParseFlags(version string) (*CLIFlags, error)
		ValidateFlags(flags *CLIFlags) error
	}

	// DataStore defines the interface to manage the data
	DataStore interface {
		Open() error
		Init() error
		Close() error
		MigrateData() error
	}

	// Server defines the interface to serve the API
	Server interface {
		Start() error
	}

	// UserService represents a service for managing user data
	UserService interface {
		User(ID UserID) (*User, error)
		UserByUsername(username string) (*User, error)
		Users() ([]User, error)
		UsersByRole(role UserRole) ([]User, error)
		CreateUser(user *User) error
		UpdateUser(ID UserID, user *User) error
		DeleteUser(ID UserID) error
	}

	// TeamService represents a service for managing user data
	TeamService interface {
		Team(ID TeamID) (*Team, error)
		TeamByName(name string) (*Team, error)
		Teams() ([]Team, error)
		CreateTeam(team *Team) error
		UpdateTeam(ID TeamID, team *Team) error
		DeleteTeam(ID TeamID) error
	}

	// TeamMembershipService represents a service for managing team membership data
	TeamMembershipService interface {
		TeamMembership(ID TeamMembershipID) (*TeamMembership, error)
		TeamMemberships() ([]TeamMembership, error)
		TeamMembershipsByUserID(userID UserID) ([]TeamMembership, error)
		TeamMembershipsByTeamID(teamID TeamID) ([]TeamMembership, error)
		CreateTeamMembership(membership *TeamMembership) error
		UpdateTeamMembership(ID TeamMembershipID, membership *TeamMembership) error
		DeleteTeamMembership(ID TeamMembershipID) error
		DeleteTeamMembershipByUserID(userID UserID) error
		DeleteTeamMembershipByTeamID(teamID TeamID) error
	}

	// Baask8sService represents a service for managing baask8s data
	Baask8sService interface {
		Baask8s(ID Baask8sID) (*Baask8s, error)
		Baask8ss() ([]Baask8s, error)
		CreateBaask8s(baask8s *Baask8s) error
		UpdateBaask8s(ID Baask8sID, baask8s *Baask8s) error
		DeleteBaask8s(ID Baask8sID) error
		GetNextIdentifier() int
	}

	// BaasmspService represents a service for managing baasmsp data
	BaasmspService interface {
		Baasmsp(ID BaasmspID) (*Baasmsp, error)
		Baasmsps() ([]Baasmsp, error)
		CreateBaasmsp(baasmsp *Baasmsp) error
		DeleteBaasmsp(ID BaasmspID) error
		GetNextIdentifier() int
	}

	// Baask8sGroupService represents a service for managing baask8s group data
	Baask8sGroupService interface {
		Baask8sGroup(ID Baask8sGroupID) (*Baask8sGroup, error)
		Baask8sGroups() ([]Baask8sGroup, error)
		CreateBaask8sGroup(group *Baask8sGroup) error
		UpdateBaask8sGroup(ID Baask8sGroupID, group *Baask8sGroup) error
		DeleteBaask8sGroup(ID Baask8sGroupID) error
	}

	// RegistryService represents a service for managing registry data
	RegistryService interface {
		Registry(ID RegistryID) (*Registry, error)
		Registries() ([]Registry, error)
		CreateRegistry(registry *Registry) error
		UpdateRegistry(ID RegistryID, registry *Registry) error
		DeleteRegistry(ID RegistryID) error
	}

	// SettingsService represents a service for managing application settings
	SettingsService interface {
		Settings() (*Settings, error)
		UpdateSettings(settings *Settings) error
	}

	// VersionService represents a service for managing version data
	VersionService interface {
		DBVersion() (int, error)
		StoreDBVersion(version int) error
	}

	// WebhookService represents a service for managing webhook data.
	WebhookService interface {
		Webhooks() ([]Webhook, error)
		Webhook(ID WebhookID) (*Webhook, error)
		CreateWebhook(baasapi *Webhook) error
		WebhookByResourceID(resourceID string) (*Webhook, error)
		WebhookByToken(token string) (*Webhook, error)
		DeleteWebhook(serviceID WebhookID) error
	}

	// ResourceControlService represents a service for managing resource control data
	ResourceControlService interface {
		ResourceControl(ID ResourceControlID) (*ResourceControl, error)
		ResourceControlByResourceID(resourceID string) (*ResourceControl, error)
		ResourceControls() ([]ResourceControl, error)
		CreateResourceControl(rc *ResourceControl) error
		UpdateResourceControl(ID ResourceControlID, resourceControl *ResourceControl) error
		DeleteResourceControl(ID ResourceControlID) error
	}

	// ScheduleService represents a service for managing schedule data
	ScheduleService interface {
		Schedule(ID ScheduleID) (*Schedule, error)
		Schedules() ([]Schedule, error)
		SchedulesByJobType(jobType JobType) ([]Schedule, error)
		CreateSchedule(schedule *Schedule) error
		UpdateSchedule(ID ScheduleID, schedule *Schedule) error
		DeleteSchedule(ID ScheduleID) error
		GetNextIdentifier() int
	}

	// TagService represents a service for managing tag data
	TagService interface {
		Tags() ([]Tag, error)
		CreateTag(tag *Tag) error
		DeleteTag(ID TagID) error
	}

	// TemplateService represents a service for managing template data
	TemplateService interface {
		Templates() ([]Template, error)
		Template(ID TemplateID) (*Template, error)
		CreateTemplate(template *Template) error
		UpdateTemplate(ID TemplateID, template *Template) error
		DeleteTemplate(ID TemplateID) error
	}

	// ExtensionService represents a service for managing extension data
	ExtensionService interface {
		Extension(ID ExtensionID) (*Extension, error)
		Extensions() ([]Extension, error)
		Persist(extension *Extension) error
		DeleteExtension(ID ExtensionID) error
	}

	// CryptoService represents a service for encrypting/hashing data
	CryptoService interface {
		Hash(data string) (string, error)
		CompareHashAndData(hash string, data string) error
	}

	// DigitalSignatureService represents a service to manage digital signatures
	DigitalSignatureService interface {
		ParseKeyPair(private, public []byte) error
		GenerateKeyPair() ([]byte, []byte, error)
		EncodedPublicKey() string
		PEMHeaders() (string, string)
		CreateSignature(message string) (string, error)
	}

	// JWTService represents a service for managing JWT tokens
	JWTService interface {
		GenerateToken(data *TokenData) (string, error)
		ParseAndVerifyToken(token string) (*TokenData, error)
	}

	// FileService represents a service for managing files
	FileService interface {
		GetFileContent(filePath string) ([]byte, error)
		Rename(oldPath, newPath string) error
		RemoveDirectory(directoryPath string) error
		StoreTLSFileFromBytes(folder string, fileType TLSFileType, data []byte) (string, error)
		StoreKubeconfigFileFromBytes(folder string, fileType TLSFileType, data []byte) (string, error)
		GetPathForTLSFile(folder string, fileType TLSFileType) (string, error)
		DeleteTLSFile(folder string, fileType TLSFileType) error
		DeleteTLSFiles(folder string) error
		GetStackProjectPath(stackIdentifier string) string
		StoreStackFileFromBytes(stackIdentifier, fileName string, data []byte) (string, error)
		StoreRegistryManagementFileFromBytes(folder, fileName string, data []byte) (string, error)
		KeyPairFilesExist() (bool, error)
		StoreKeyPair(private, public []byte, privatePEMHeader, publicPEMHeader string) error
		LoadKeyPair() ([]byte, []byte, error)
		WriteJSONToFile(path string, content interface{}) error
		StoreYamlFileFromJSON(stackIdentifier, fileName string, content interface{}, creator string) (string, error) 
		FileExists(path string) (bool, error)
		StoreScheduledJobFileFromBytes(identifier string, data []byte) (string, error)
		GetScheduleFolder(identifier string) string
		ExtractExtensionArchive(data []byte) error
		GetBinaryFolder() string
	}

	// GitService represents a service for managing Git
	GitService interface {
		ClonePublicRepository(repositoryURL, referenceName string, destination string) error
		ClonePrivateRepositoryWithBasicAuth(repositoryURL, referenceName string, destination, username, password string) error
	}

	// JobScheduler represents a service to run jobs on a periodic basis
	JobScheduler interface {
		ScheduleJob(runner JobRunner) error
		UpdateJobSchedule(runner JobRunner) error
		UpdateSystemJobSchedule(jobType JobType, newCronExpression string) error
		UnscheduleJob(ID ScheduleID)
		Start()
	}

	// JobRunner represents a service that can be used to run a job
	JobRunner interface {
		Run()
		GetSchedule() *Schedule
	}

	// Snapshotter represents a service used to create baask8s snapshots
	Snapshotter interface {
		CreateSnapshot(baask8s *Baask8s) (*Snapshot, error)
	}

	// LDAPService represents a service used to authenticate users against a LDAP/AD
	LDAPService interface {
		AuthenticateUser(username, password string, settings *LDAPSettings) error
		TestConnectivity(settings *LDAPSettings) error
		GetUserGroups(username string, settings *LDAPSettings) ([]string, error)
	}

	// CAFilesManager represents a service to manage CAFiles 
	CAFilesManager interface {
		Deploy(creator string, networkID string, ansible_extra string, ansible_env string, ansible_config string, flag bool) error
		GetLogs(namespade string, nline int) (string, error)
		//Remove(baask8s *Baask8s, creator string, networkID string) error
	}

	// JobService represents a service to manage job execution on hosts
	JobService interface {
		ExecuteScript(baask8s *Baask8s, nodeName, image string, script []byte, schedule *Schedule) error
	}

	// ExtensionManager represents a service used to manage extensions
	ExtensionManager interface {
		FetchExtensionDefinitions() ([]Extension, error)
		EnableExtension(extension *Extension, licenseKey string) error
		DisableExtension(extension *Extension) error
		UpdateExtension(extension *Extension, version string) error
	}
)

const (
	// APIVersion is the version number of the BaaSapi API
	APIVersion = "1.20.3"
	// DBVersion is the version number of the BaaSapi database
	DBVersion = 18
	// AssetsServerURL represents the URL of the BaaSapi asset server
	AssetsServerURL = "https://baasapi-io-assets.sfo2.digitaloceanspaces.com"
	// MessageOfTheDayURL represents the URL where BaaSapi MOTD message can be retrieved
	MessageOfTheDayURL = AssetsServerURL + "/motd.html"
	// MessageOfTheDayTitleURL represents the URL where BaaSapi MOTD title can be retrieved
	MessageOfTheDayTitleURL = AssetsServerURL + "/motd-title.txt"
	// ExtensionDefinitionsURL represents the URL where BaaSapi extension definitions can be retrieved
	ExtensionDefinitionsURL = AssetsServerURL + "/extensions-1.20.2.json"
	// BaaSapiAgentHeader represents the name of the header available in any agent response
	BaaSapiAgentHeader = "BaaSapi-Agent"
	// BaaSapiAgentTargetHeader represent the name of the header containing the target node name
	BaaSapiAgentTargetHeader = "X-BaaSapiAgent-Target"
	// BaaSapiAgentSignatureHeader represent the name of the header containing the digital signature
	BaaSapiAgentSignatureHeader = "X-BaaSapiAgent-Signature"
	// BaaSapiAgentPublicKeyHeader represent the name of the header containing the public key
	BaaSapiAgentPublicKeyHeader = "X-BaaSapiAgent-PublicKey"
	// BaaSapiAgentSignatureMessage represents the message used to create a digital signature
	// to be used when communicating with an agent
	BaaSapiAgentSignatureMessage = "BaaSapi-App"
	// SupportedDockerAPIVersion is the minimum Docker API version supported by BaaSapi
	SupportedDockerAPIVersion = "1.24"
)

const (
	// TLSFileCA represents a TLS CA certificate file
	TLSFileCA TLSFileType = iota
	// TLSFileCert represents a TLS certificate file
	TLSFileCert
	// TLSFileKey represents a TLS key file
	TLSFileKey
	// Other file
	TLSFileOther
)

const (
	_ MembershipRole = iota
	// TeamLeader represents a leader role inside a team
	TeamLeader
	// TeamMember represents a member role inside a team
	TeamMember
)

const (
	_ UserRole = iota
	// AdministratorRole represents an administrator user role
	AdministratorRole
	// StandardUserRole represents a regular user role
	StandardUserRole
)

const (
	_ AuthenticationMethod = iota
	// AuthenticationInternal represents the internal authentication method (authentication against BaaSapi API)
	AuthenticationInternal
	// AuthenticationLDAP represents the LDAP authentication method (authentication against a LDAP server)
	AuthenticationLDAP
	//AuthenticationOAuth represents the OAuth authentication method (authentication against a authorization server)
	AuthenticationOAuth
)

const (
	_ ResourceAccessLevel = iota
	// ReadWriteAccessLevel represents an access level with read-write permissions on a resource
	ReadWriteAccessLevel
)

const (
	_ ResourceControlType = iota
	// ContainerResourceControl represents a resource control associated to a Docker container
	ContainerResourceControl
	// ServiceResourceControl represents a resource control associated to a Docker service
	ServiceResourceControl
	// VolumeResourceControl represents a resource control associated to a Docker volume
	VolumeResourceControl
	// NetworkResourceControl represents a resource control associated to a Docker network
	NetworkResourceControl
	// SecretResourceControl represents a resource control associated to a Docker secret
	SecretResourceControl
	// StackResourceControl represents a resource control associated to a stack composed of Docker services
	StackResourceControl
	// ConfigResourceControl represents a resource control associated to a Docker config
	ConfigResourceControl
)

const (
	_ Baask8sType = iota
	// DockerEnvironment represents an baask8s connected to a Docker environment
	DockerEnvironment2
	// AgentOnDockerEnvironment represents an baask8s connected to a BaaSapi agent deployed on a Docker environment
	AgentOnDockerEnvironment2
	// AzureEnvironment represents an baask8s connected to an Azure environment
	AzureEnvironment2
)

const (
	_ TemplateType = iota
	// ContainerTemplate represents a container template
	ContainerTemplate
	// SwarmStackTemplate represents a template used to deploy a Swarm stack
	SwarmStackTemplate
	// ComposeStackTemplate represents a template used to deploy a Compose stack
	ComposeStackTemplate
)

const (
	_ Baask8sStatus = iota
	// Baask8sStatusUp is used to represent an available baask8s
	Baask8sStatusUp
	// Baask8sStatusDown is used to represent an unavailable baask8s
	Baask8sStatusDown
)

const (
	_ WebhookType = iota
	// ServiceWebhook is a webhook for restarting a docker service
	ServiceWebhook
)

const (
	_ ExtensionID = iota
	// RegistryManagementExtension represents the registry management extension
	RegistryManagementExtension
	// OAuthAuthenticationExtension represents the OAuth authentication extension
	OAuthAuthenticationExtension
)

const (
	_ JobType = iota
	// ScriptExecutionJobType is a non-system job used to execute a script against a list of
	// baask8ss via privileged containers
	ScriptExecutionJobType
	// SnapshotJobType is a system job used to create baask8s snapshots
	SnapshotJobType
	// Baask8sJobType is a system job used to sync baask8s
	Baask8sJobType
	Baask8sSyncJobType
)

const (
	_ RegistryType = iota
	// QuayRegistry represents a Quay.io registry
	QuayRegistry
	// AzureRegistry represents an ACR registry
	AzureRegistry
	// CustomRegistry represents a custom registry
	CustomRegistry
)
