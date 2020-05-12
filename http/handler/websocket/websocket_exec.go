package websocket

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	httperror "github.com/baasapi/libhttp/error"
	"github.com/baasapi/libhttp/request"
	"github.com/baasapi/baasapi/api"
	"github.com/baasapi/baasapi/api/crypto"
)

type webSocketExecRequestParams struct {
	execID   string
	nodeName string
	baask8s *baasapi.Baask8s
}

type execStartOperationPayload struct {
	Tty    bool
	Detach bool
}

// websocketExec handles GET requests on /websocket/exec?id=<execID>&baask8sId=<baask8sID>&nodeName=<nodeName>&token=<token>
// If the nodeName query parameter is present, the request will be proxied to the underlying agent baask8s.
// If the nodeName query parameter is not specified, the request will be upgraded to the websocket protocol and
// an ExecStart operation HTTP request will be created and hijacked.
// Authentication and access is controled via the mandatory token query parameter.
func (handler *Handler) websocketExec(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	execID, err := request.RetrieveQueryParameter(r, "id", false)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: id", err}
	}
	if !govalidator.IsHexadecimal(execID) {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: id (must be hexadecimal identifier)", err}
	}

	baask8sID, err := request.RetrieveNumericQueryParameter(r, "baask8sId", false)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid query parameter: baask8sId", err}
	}

	baask8s, err := handler.Baask8sService.Baask8s(baasapi.Baask8sID(baask8sID))
	if err == baasapi.ErrObjectNotFound {
		return &httperror.HandlerError{http.StatusNotFound, "Unable to find the baask8s associated to the stack inside the database", err}
	} else if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "Unable to find the baask8s associated to the stack inside the database", err}
	}

	err = handler.requestBouncer.Baask8sAccess(r, baask8s)
	if err != nil {
		return &httperror.HandlerError{http.StatusForbidden, "Permission denied to access baask8s", baasapi.ErrBaask8sAccessDenied}
	}

	params := &webSocketExecRequestParams{
		baask8s: baask8s,
		execID:   execID,
		nodeName: r.FormValue("nodeName"),
	}

	err = handler.handleRequest(w, r, params)
	if err != nil {
		return &httperror.HandlerError{http.StatusInternalServerError, "An error occured during websocket exec operation", err}
	}

	return nil
}

func (handler *Handler) handleRequest(w http.ResponseWriter, r *http.Request, params *webSocketExecRequestParams) error {
	r.Header.Del("Origin")

	if params.nodeName != "" || params.baask8s.Type == baasapi.AgentOnDockerEnvironment {
		return handler.proxyWebsocketRequest(w, r, params)
	}

	websocketConn, err := handler.connectionUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer websocketConn.Close()

	return hijackExecStartOperation(websocketConn, params.baask8s, params.execID)
}

func (handler *Handler) proxyWebsocketRequest(w http.ResponseWriter, r *http.Request, params *webSocketExecRequestParams) error {
	agentURL, err := url.Parse(params.baask8s.URL)
	if err != nil {
		return err
	}

	agentURL.Scheme = "ws"
	proxy := websocketproxy.NewProxy(agentURL)

	//if params.baask8s.TLSConfig.TLS || params.baask8s.TLSConfig.TLSSkipVerify {
	//	agentURL.Scheme = "wss"
	//	proxy.Dialer = &websocket.Dialer{
	//		TLSClientConfig: &tls.Config{
	//			InsecureSkipVerify: params.baask8s.TLSConfig.TLSSkipVerify,
	//		},
	//	}
	//}

	signature, err := handler.SignatureService.CreateSignature(baasapi.BaaSapiAgentSignatureMessage)
	if err != nil {
		return err
	}

	proxy.Director = func(incoming *http.Request, out http.Header) {
		out.Set(baasapi.BaaSapiAgentPublicKeyHeader, handler.SignatureService.EncodedPublicKey())
		out.Set(baasapi.BaaSapiAgentSignatureHeader, signature)
		out.Set(baasapi.BaaSapiAgentTargetHeader, params.nodeName)
	}

	proxy.ServeHTTP(w, r)

	return nil
}

func hijackExecStartOperation(websocketConn *websocket.Conn, baask8s *baasapi.Baask8s, execID string) error {
	dial, err := initDial(baask8s)
	if err != nil {
		return err
	}

	// When we set up a TCP connection for hijack, there could be long periods
	// of inactivity (a long running command with no output) that in certain
	// network setups may cause ECONNTIMEOUT, leaving the client in an unknown
	// state. Setting TCP KeepAlive on the socket connection will prohibit
	// ECONNTIMEOUT unless the socket connection truly is broken
	if tcpConn, ok := dial.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	httpConn := httputil.NewClientConn(dial, nil)
	defer httpConn.Close()

	execStartRequest, err := createExecStartRequest(execID)
	if err != nil {
		return err
	}

	err = hijackRequest(websocketConn, httpConn, execStartRequest)
	if err != nil {
		return err
	}

	return nil
}

func initDial(baask8s *baasapi.Baask8s) (net.Conn, error) {
	url, err := url.Parse(baask8s.URL)
	if err != nil {
		return nil, err
	}

	host := url.Host

	if url.Scheme == "unix" || url.Scheme == "npipe" {
		host = url.Path
	}

	//if baask8s.TLSConfig.TLS {
	//	tlsConfig, err := crypto.CreateTLSConfigurationFromDisk(baask8s.TLSConfig.TLSCACertPath, baask8s.TLSConfig.TLSCertPath, baask8s.TLSConfig.TLSKeyPath, baask8s.TLSConfig.TLSSkipVerify)
	//	if err != nil {
	//		return nil, err
	//	}

	//	return tls.Dial(url.Scheme, host, tlsConfig)
	//}

	return createDial(url.Scheme, host)
}

func createExecStartRequest(execID string) (*http.Request, error) {
	execStartOperationPayload := &execStartOperationPayload{
		Tty:    true,
		Detach: false,
	}

	encodedBody := bytes.NewBuffer(nil)
	err := json.NewEncoder(encodedBody).Encode(execStartOperationPayload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", "/exec/"+execID+"/start", encodedBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "tcp")

	return request, nil
}

func hijackRequest(websocketConn *websocket.Conn, httpConn *httputil.ClientConn, request *http.Request) error {
	// Server hijacks the connection, error 'connection closed' expected
	resp, err := httpConn.Do(request)
	if err != httputil.ErrPersistEOF {
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusSwitchingProtocols {
			resp.Body.Close()
			return fmt.Errorf("unable to upgrade to tcp, received %d", resp.StatusCode)
		}
	}

	tcpConn, brw := httpConn.Hijack()
	defer tcpConn.Close()

	errorChan := make(chan error, 1)
	go streamFromTCPConnToWebsocketConn(websocketConn, brw, errorChan)
	go streamFromWebsocketConnToTCPConn(websocketConn, tcpConn, errorChan)

	err = <-errorChan
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
		return err
	}

	return nil
}

func streamFromWebsocketConnToTCPConn(websocketConn *websocket.Conn, tcpConn net.Conn, errorChan chan error) {
	for {
		_, in, err := websocketConn.ReadMessage()
		if err != nil {
			errorChan <- err
			break
		}

		_, err = tcpConn.Write(in)
		if err != nil {
			errorChan <- err
			break
		}
	}
}

func streamFromTCPConnToWebsocketConn(websocketConn *websocket.Conn, br *bufio.Reader, errorChan chan error) {
	for {
		out := make([]byte, 2048)
		_, err := br.Read(out)
		if err != nil {
			errorChan <- err
			break
		}

		err = websocketConn.WriteMessage(websocket.TextMessage, out)
		if err != nil {
			errorChan <- err
			break
		}
	}
}
