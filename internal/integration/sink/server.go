//go:build integration

package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"sync"
	"sync/atomic"
	"time"

	crewjam_saml "github.com/crewjam/saml"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/idp/providers/saml"
)

const (
	port       = "8081"
	listenAddr = "127.0.0.1:" + port
	host       = "localhost:" + port
)

// CallURL returns the full URL to the handler of a [Channel].
func CallURL(ch Channel) string {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   rootPath(ch),
	}
	return u.String()
}

func SuccessfulOAuthIntent(instanceID, idpID, idpUserID, userID string) (string, string, time.Time, uint64, error) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   successfulIntentOAuthPath(),
	}
	resp, err := callIntent(u.String(), &SuccessfulIntentRequest{
		InstanceID: instanceID,
		IDPID:      idpID,
		IDPUserID:  idpUserID,
		UserID:     userID,
	})
	if err != nil {
		return "", "", time.Time{}, uint64(0), err
	}
	return resp.IntentID, resp.Token, resp.ChangeDate, resp.Sequence, nil
}

func SuccessfulOIDCIntent(instanceID, idpID, idpUserID, userID string) (string, string, time.Time, uint64, error) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   successfulIntentOIDCPath(),
	}
	resp, err := callIntent(u.String(), &SuccessfulIntentRequest{
		InstanceID: instanceID,
		IDPID:      idpID,
		IDPUserID:  idpUserID,
		UserID:     userID,
	})
	if err != nil {
		return "", "", time.Time{}, uint64(0), err
	}
	return resp.IntentID, resp.Token, resp.ChangeDate, resp.Sequence, nil
}

func SuccessfulSAMLIntent(instanceID, idpID, idpUserID, userID string) (string, string, time.Time, uint64, error) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   successfulIntentSAMLPath(),
	}
	resp, err := callIntent(u.String(), &SuccessfulIntentRequest{
		InstanceID: instanceID,
		IDPID:      idpID,
		IDPUserID:  idpUserID,
		UserID:     userID,
	})
	if err != nil {
		return "", "", time.Time{}, uint64(0), err
	}
	return resp.IntentID, resp.Token, resp.ChangeDate, resp.Sequence, nil
}

func SuccessfulLDAPIntent(instanceID, idpID, idpUserID, userID string) (string, string, time.Time, uint64, error) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   successfulIntentLDAPPath(),
	}
	resp, err := callIntent(u.String(), &SuccessfulIntentRequest{
		InstanceID: instanceID,
		IDPID:      idpID,
		IDPUserID:  idpUserID,
		UserID:     userID,
	})
	if err != nil {
		return "", "", time.Time{}, uint64(0), err
	}
	return resp.IntentID, resp.Token, resp.ChangeDate, resp.Sequence, nil
}

// StartServer starts a simple HTTP server on localhost:8081
// ZITADEL can use the server to send HTTP requests which can be
// used to validate tests through [Subscribe]rs.
// For each [Channel] a route is registered on http://localhost:8081/<channel_name>.
// The route must be used to send the HTTP request to be validated.
// [CallURL] can be used to obtain the full URL for a given Channel.
//
// This function is only active when the `integration` build tag is enabled
func StartServer(commands *command.Commands) (close func()) {
	router := chi.NewRouter()
	for _, ch := range ChannelValues() {
		fwd := &forwarder{
			channelID:   ch,
			subscribers: make(map[int64]chan<- *Request),
		}
		router.HandleFunc(rootPath(ch), fwd.receiveHandler)
		router.HandleFunc(subscribePath(ch), fwd.subscriptionHandler)
		router.HandleFunc(successfulIntentOAuthPath(), successfulIntentHandler(commands, createSuccessfulOAuthIntent))
		router.HandleFunc(successfulIntentOIDCPath(), successfulIntentHandler(commands, createSuccessfulOIDCIntent))
		router.HandleFunc(successfulIntentSAMLPath(), successfulIntentHandler(commands, createSuccessfulSAMLIntent))
		router.HandleFunc(successfulIntentLDAPPath(), successfulIntentHandler(commands, createSuccessfulLDAPIntent))
	}
	s := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	logging.WithFields("listen_addr", listenAddr).Warn("!!!! A sink server is started which may expose sensitive data on a public endpoint. Make sure the `integration` build tag is disabled for production builds. !!!!")
	go func() {
		err := s.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logging.WithError(err).Fatal("sink server")
		}
	}()
	return func() {
		logging.OnError(s.Close()).Error("sink server")
	}
}

func rootPath(c Channel) string {
	return path.Join("/", c.String())
}

func subscribePath(c Channel) string {
	return path.Join("/", c.String(), "subscribe")
}

func intentPath() string {
	return path.Join("/", "intent")
}

func successfulIntentPath() string {
	return path.Join(intentPath(), "/", "successful")
}

func successfulIntentOAuthPath() string {
	return path.Join(successfulIntentPath(), "/", "oauth")
}

func successfulIntentOIDCPath() string {
	return path.Join(successfulIntentPath(), "/", "oidc")
}

func successfulIntentSAMLPath() string {
	return path.Join(successfulIntentPath(), "/", "saml")
}

func successfulIntentLDAPPath() string {
	return path.Join(successfulIntentPath(), "/", "ldap")
}

// forwarder handles incoming HTTP requests from ZITADEL and
// forwards them to all subscribed web sockets.
type forwarder struct {
	channelID   Channel
	id          atomic.Int64
	mtx         sync.RWMutex
	subscribers map[int64]chan<- *Request
	upgrader    websocket.Upgrader
}

// receiveHandler receives a simple HTTP for a single [Channel]
// and forwards them on all active subscribers of that Channel.
func (c *forwarder) receiveHandler(w http.ResponseWriter, r *http.Request) {
	req := &Request{
		Header: r.Header.Clone(),
	}
	var err error
	req.Body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c.mtx.RLock()
	for _, reqChan := range c.subscribers {
		reqChan <- req
	}
	c.mtx.RUnlock()
	w.WriteHeader(http.StatusOK)
}

// subscriptionHandler upgrades HTTP request to a websocket connection for subscribers.
// All received HTTP requests on a subscriber's channel are send on the websocket to the client.
func (c *forwarder) subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := c.upgrader.Upgrade(w, r, nil)
	logging.OnError(err).Error("websocket upgrade")
	if err != nil {
		return
	}
	done := readLoop(ws)

	id := c.id.Add(1)
	reqChannel := make(chan *Request, 100)

	c.mtx.Lock()
	c.subscribers[id] = reqChannel
	c.mtx.Unlock()

	logging.WithFields("id", id, "channel", c.channelID).Info("websocket opened")

	defer func() {
		c.mtx.Lock()
		delete(c.subscribers, id)
		c.mtx.Unlock()

		ws.Close()
		close(reqChannel)
	}()

	for {
		select {
		case err := <-done:
			logging.WithError(err).WithFields(logrus.Fields{"id": id, "channel": c.channelID}).Info("websocket closed")
			return
		case req := <-reqChannel:
			if err := ws.WriteJSON(req); err != nil {
				logging.WithError(err).WithFields(logrus.Fields{"id": id, "channel": c.channelID}).Error("websocket write json")
				return
			}
		}
	}
}

// readLoop makes sure we can receive close messages
func readLoop(ws *websocket.Conn) (done chan error) {
	done = make(chan error, 1)

	go func(done chan<- error) {
		for {
			_, _, err := ws.NextReader()
			if err != nil {
				done <- err
				break
			}
		}
		close(done)
	}(done)

	return done
}

type SuccessfulIntentRequest struct {
	InstanceID string `json:"instance_id"`
	IDPID      string `json:"idp_id"`
	IDPUserID  string `json:"idp_user_id"`
	UserID     string `json:"user_id"`
}
type SuccessfulIntentResponse struct {
	IntentID   string    `json:"intent_id"`
	Token      string    `json:"token"`
	ChangeDate time.Time `json:"change_date"`
	Sequence   uint64    `json:"sequence"`
}

func callIntent(url string, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	result := new(SuccessfulIntentResponse)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}

func successfulIntentHandler(cmd *command.Commands, createIntent func(ctx context.Context, cmd *command.Commands, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &SuccessfulIntentRequest{}
		if err := json.Unmarshal(body, req); err != nil {
		}

		ctx := authz.WithInstanceID(r.Context(), req.InstanceID)
		resp, err := createIntent(ctx, cmd, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(data)
		return
	}
}

func createIntent(ctx context.Context, cmd *command.Commands, instanceID, idpID string) (string, error) {
	writeModel, _, err := cmd.CreateIntent(ctx, "", idpID, "https://example.com/success", "https://example.com/failure", instanceID, nil)
	if err != nil {
		return "", err
	}
	return writeModel.AggregateID, nil
}

func createSuccessfulOAuthIntent(ctx context.Context, cmd *command.Commands, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error) {
	intentID, err := createIntent(ctx, cmd, req.InstanceID, req.IDPID)
	if err != nil {
		return nil, err
	}
	writeModel, err := cmd.GetIntentWriteModel(ctx, intentID, req.InstanceID)
	if err != nil {
		return nil, err
	}
	idAttribute := "id"
	idpUser := oauth.NewUserMapper(idAttribute)
	idpUser.RawInfo = map[string]interface{}{
		idAttribute:          req.IDPUserID,
		"preferred_username": "username",
	}
	idpSession := &oauth.Session{
		Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
			Token: &oauth2.Token{
				AccessToken: "accessToken",
			},
			IDToken: "idToken",
		},
	}
	token, err := cmd.SucceedIDPIntent(ctx, writeModel, idpUser, idpSession, req.UserID)
	if err != nil {
		return nil, err
	}
	return &SuccessfulIntentResponse{
		intentID,
		token,
		writeModel.ChangeDate,
		writeModel.ProcessedSequence,
	}, nil
}

func createSuccessfulOIDCIntent(ctx context.Context, cmd *command.Commands, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error) {
	intentID, err := createIntent(ctx, cmd, req.InstanceID, req.IDPID)
	writeModel, err := cmd.GetIntentWriteModel(ctx, intentID, req.InstanceID)
	idpUser := openid.NewUser(
		&oidc.UserInfo{
			Subject: req.IDPUserID,
			UserInfoProfile: oidc.UserInfoProfile{
				PreferredUsername: "username",
			},
		},
	)
	idpSession := &openid.Session{
		Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
			Token: &oauth2.Token{
				AccessToken: "accessToken",
			},
			IDToken: "idToken",
		},
	}
	token, err := cmd.SucceedIDPIntent(ctx, writeModel, idpUser, idpSession, req.UserID)
	if err != nil {
		return nil, err
	}
	return &SuccessfulIntentResponse{
		intentID,
		token,
		writeModel.ChangeDate,
		writeModel.ProcessedSequence,
	}, nil
}

func createSuccessfulSAMLIntent(ctx context.Context, cmd *command.Commands, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error) {
	intentID, err := createIntent(ctx, cmd, req.InstanceID, req.IDPID)
	writeModel, err := cmd.GetIntentWriteModel(ctx, intentID, req.InstanceID)

	idpUser := &saml.UserMapper{
		ID:         req.IDPUserID,
		Attributes: map[string][]string{"attribute1": {"value1"}},
	}
	assertion := &crewjam_saml.Assertion{ID: "id"}

	token, err := cmd.SucceedSAMLIDPIntent(ctx, writeModel, idpUser, req.UserID, assertion)
	if err != nil {
		return nil, err
	}
	return &SuccessfulIntentResponse{
		intentID,
		token,
		writeModel.ChangeDate,
		writeModel.ProcessedSequence,
	}, nil
}

func createSuccessfulLDAPIntent(ctx context.Context, cmd *command.Commands, req *SuccessfulIntentRequest) (*SuccessfulIntentResponse, error) {
	intentID, err := createIntent(ctx, cmd, req.InstanceID, req.IDPID)
	writeModel, err := cmd.GetIntentWriteModel(ctx, intentID, req.InstanceID)
	username := "username"
	lang := language.Make("en")
	idpUser := ldap.NewUser(
		req.IDPUserID,
		"",
		"",
		"",
		"",
		username,
		"",
		false,
		"",
		false,
		lang,
		"",
		"",
	)
	attributes := map[string][]string{"id": {req.IDPUserID}, "username": {username}, "language": {lang.String()}}
	token, err := cmd.SucceedLDAPIDPIntent(ctx, writeModel, idpUser, req.UserID, attributes)
	if err != nil {
		return nil, err
	}
	return &SuccessfulIntentResponse{
		intentID,
		token,
		writeModel.ChangeDate,
		writeModel.ProcessedSequence,
	}, nil
}
