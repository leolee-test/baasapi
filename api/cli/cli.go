package cli

import (
	"time"

	"github.com/baasapi/baasapi/api"

	"os"
	"path/filepath"
	//"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Service implements the CLIService interface
type Service struct{}

const (
	errSocketOrNamedPipeNotFound     = baasapi.Error("Unable to locate Unix socket or named pipe")
	errTemplateFileNotFound          = baasapi.Error("Unable to locate template file on disk")
	errInvalidSyncInterval           = baasapi.Error("Invalid synchronization interval")
	errInvalidSnapshotInterval       = baasapi.Error("Invalid snapshot interval")
	errInvalidBaask8sInterval        = baasapi.Error("Invalid baask8s interval")
	errNoAuthExcludeAdminPassword    = baasapi.Error("Cannot use --no-auth with --admin-password or --admin-password-file")
	errAdminPassExcludeAdminPassFile = baasapi.Error("Cannot use --admin-password with --admin-password-file")
)

// ParseFlags parse the CLI flags and return a baasapi.Flags struct
func (*Service) ParseFlags(version string) (*baasapi.CLIFlags, error) {
	kingpin.Version(version)

	flags := &baasapi.CLIFlags{
		Addr:              kingpin.Flag("bind", "Address and port to serve BaaSapi").Default(defaultBindAddress).Short('p').String(),
		Assets:            kingpin.Flag("assets", "Path to the assets").Default(defaultAssetsDirectory).Short('a').String(),
		Data:              kingpin.Flag("data", "Path to the folder where the data is stored").Default(defaultDataDirectory).Short('d').String(),
		NoAuth:            kingpin.Flag("no-auth", "Disable authentication").Default(defaultNoAuth).Bool(),
		NoAnalytics:       kingpin.Flag("no-analytics", "Disable Analytics in app").Default(defaultNoAnalytics).Bool(),
		TLS:               kingpin.Flag("tlsverify", "TLS support").Default(defaultTLS).Bool(),
		TLSSkipVerify:     kingpin.Flag("tlsskipverify", "Disable TLS server verification").Default(defaultTLSSkipVerify).Bool(),
		TLSCacert:         kingpin.Flag("tlscacert", "Path to the CA").Default(defaultTLSCACertPath).String(),
		TLSCert:           kingpin.Flag("tlscert", "Path to the TLS certificate file").Default(defaultTLSCertPath).String(),
		TLSKey:            kingpin.Flag("tlskey", "Path to the TLS key").Default(defaultTLSKeyPath).String(),
		SSL:               kingpin.Flag("ssl", "Secure BaaSapi instance using SSL").Default(defaultSSL).Bool(),
		SSLCert:           kingpin.Flag("sslcert", "Path to the SSL certificate used to secure the BaaSapi instance").Default(defaultSSLCertPath).String(),
		SSLKey:            kingpin.Flag("sslkey", "Path to the SSL key used to secure the BaaSapi instance").Default(defaultSSLKeyPath).String(),
		SyncInterval:      kingpin.Flag("sync-interval", "Duration between each synchronization via the external baask8ss source").Default(defaultSyncInterval).String(),
		Snapshot:          kingpin.Flag("snapshot", "Start a background job to create baask8s snapshots").Default(defaultSnapshot).Bool(),
		SnapshotInterval:  kingpin.Flag("snapshot-interval", "Duration between each baask8s snapshot job").Default(defaultSnapshotInterval).String(),
		Baask8s:           kingpin.Flag("baask8s", "Start a background job to sync baask8s").Default(defaultBaask8s).Bool(),
		Baask8sInterval:   kingpin.Flag("baask8s-interval", "Duration between each baask8s job").Default(defaultBaask8sInterval).String(),
		AdminPassword:     kingpin.Flag("admin-password", "Hashed admin password").String(),
		AdminPasswordFile: kingpin.Flag("admin-password-file", "Path to the file containing the password for the admin user").String(),
		Labels:            pairs(kingpin.Flag("hide-label", "Hide containers with a specific label in the UI").Short('l')),
		Logo:              kingpin.Flag("logo", "URL for the logo displayed in the UI").String(),
		Templates:         kingpin.Flag("templates", "URL to the templates definitions.").Short('t').String(),
		TemplateFile:      kingpin.Flag("template-file", "Path to the templates (app) definitions on the filesystem").Default(defaultTemplateFile).String(),
	}

	kingpin.Parse()

	if !filepath.IsAbs(*flags.Assets) {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		*flags.Assets = filepath.Join(filepath.Dir(ex), *flags.Assets)
	}

	return flags, nil
}

// ValidateFlags validates the values of the flags.
func (*Service) ValidateFlags(flags *baasapi.CLIFlags) error {

	err := validateTemplateFile(*flags.TemplateFile)
	if err != nil {
		return err
	}

	err = validateSyncInterval(*flags.SyncInterval)
	if err != nil {
		return err
	}

	err = validateSnapshotInterval(*flags.SnapshotInterval)
	if err != nil {
		return err
	}

	err = validateBaask8sInterval(*flags.Baask8sInterval)
	if err != nil {
		return err
	}

	if *flags.NoAuth && (*flags.AdminPassword != "" || *flags.AdminPasswordFile != "") {
		return errNoAuthExcludeAdminPassword
	}

	if *flags.AdminPassword != "" && *flags.AdminPasswordFile != "" {
		return errAdminPassExcludeAdminPassFile
	}

	return nil
}

func validateTemplateFile(templateFile string) error {
	if _, err := os.Stat(templateFile); err != nil {
		if os.IsNotExist(err) {
			return errTemplateFileNotFound
		}
		return err
	}
	return nil
}

func validateSyncInterval(syncInterval string) error {
	if syncInterval != defaultSyncInterval {
		_, err := time.ParseDuration(syncInterval)
		if err != nil {
			return errInvalidSyncInterval
		}
	}
	return nil
}

func validateSnapshotInterval(snapshotInterval string) error {
	if snapshotInterval != defaultSnapshotInterval {
		_, err := time.ParseDuration(snapshotInterval)
		if err != nil {
			return errInvalidSnapshotInterval
		}
	}
	return nil
}

func validateBaask8sInterval(baask8sInterval string) error {
	if baask8sInterval != defaultBaask8sInterval {
		_, err := time.ParseDuration(baask8sInterval)
		if err != nil {
			return errInvalidBaask8sInterval
		}
	}
	return nil
}
