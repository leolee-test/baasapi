// +build !windows

package cli

const (
	defaultBindAddress      = ":9000"
	defaultDataDirectory    = "/data"
	defaultAssetsDirectory  = "./"
	defaultNoAuth           = "false"
	defaultNoAnalytics      = "false"
	defaultTLS              = "false"
	defaultTLSSkipVerify    = "false"
	defaultTLSCACertPath    = "/certs/ca.pem"
	defaultTLSCertPath      = "/certs/cert.pem"
	defaultTLSKeyPath       = "/certs/key.pem"
	defaultSSL              = "false"
	defaultSSLCertPath      = "/certs/baasapi.crt"
	defaultSSLKeyPath       = "/certs/baasapi.key"
	defaultSyncInterval     = "60s"
	defaultSnapshot         = "true"
	defaultSnapshotInterval = "5m"
	defaultBaask8s          = "true"
	defaultBaask8sInterval  = "1m"
	defaultTemplateFile     = "/templates.json"
)
