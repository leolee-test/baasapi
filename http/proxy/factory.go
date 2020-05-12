package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/crypto"
)

// AzureAPIBaseURL is the URL where Azure API requests will be proxied.
const AzureAPIBaseURL = "https://management.azure.com"

// proxyFactory is a factory to create reverse proxies to Docker endpoints
type proxyFactory struct {
	ResourceControlService baasapi.ResourceControlService
	TeamMembershipService  baasapi.TeamMembershipService
	SettingsService        baasapi.SettingsService
	RegistryService        baasapi.RegistryService
	DockerHubService       baasapi.DockerHubService
	SignatureService       baasapi.DigitalSignatureService
}

func (factory *proxyFactory) newHTTPProxy(u *url.URL) http.Handler {
	u.Scheme = "http"
	return httputil.NewSingleHostReverseProxy(u)
}

func newAzureProxy(credentials *baasapi.AzureCredentials) (http.Handler, error) {
	url, err := url.Parse(AzureAPIBaseURL)
	if err != nil {
		return nil, err
	}

	proxy := newSingleHostReverseProxyWithHostHeader(url)
	proxy.Transport = NewAzureTransport(credentials)

	return proxy, nil
}

func (factory *proxyFactory) newDockerHTTPSProxy(u *url.URL, tlsConfig *baasapi.TLSConfiguration, enableSignature bool) (http.Handler, error) {
	u.Scheme = "https"

	proxy := factory.createDockerReverseProxy(u, enableSignature)
	config, err := crypto.CreateTLSConfigurationFromDisk(tlsConfig.TLSCACertPath, tlsConfig.TLSCertPath, tlsConfig.TLSKeyPath, tlsConfig.TLSSkipVerify)
	if err != nil {
		return nil, err
	}

	proxy.Transport.(*proxyTransport).dockerTransport.TLSClientConfig = config
	return proxy, nil
}

func (factory *proxyFactory) newDockerHTTPProxy(u *url.URL, enableSignature bool) http.Handler {
	u.Scheme = "http"
	return factory.createDockerReverseProxy(u, enableSignature)
}

func (factory *proxyFactory) createDockerReverseProxy(u *url.URL, enableSignature bool) *httputil.ReverseProxy {
	proxy := newSingleHostReverseProxyWithHostHeader(u)
	transport := &proxyTransport{
		enableSignature:        enableSignature,
		ResourceControlService: factory.ResourceControlService,
		TeamMembershipService:  factory.TeamMembershipService,
		SettingsService:        factory.SettingsService,
		RegistryService:        factory.RegistryService,
		DockerHubService:       factory.DockerHubService,
		dockerTransport:        &http.Transport{},
	}

	if enableSignature {
		transport.SignatureService = factory.SignatureService
	}

	proxy.Transport = transport
	return proxy
}

func newSocketTransport(socketPath string) *http.Transport {
	return &http.Transport{
		Dial: func(proto, addr string) (conn net.Conn, err error) {
			return net.Dial("unix", socketPath)
		},
	}
}
