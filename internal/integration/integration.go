// Package integration provides helpers for integration testing.
package integration

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/zitadel/oidc/v3/pkg/client"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/yaml"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/net"
	"github.com/zitadel/zitadel/internal/webauthn"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/instance"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/org"
	"github.com/zitadel/zitadel/pkg/grpc/system"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

var (
	//go:embed config/client.yaml
	clientYAML []byte
	//go:embed config/system-user-key.pem
	systemUserKey []byte
)

var tmpDir string

func init() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	tmpDir = filepath.Join(string(bytes.TrimSpace(out)), "tmp")
}

// NotEmpty can be used as placeholder, when the returned values is unknown.
// It can be used in tests to assert whether a value should be empty or not.
const NotEmpty = "not empty"

const (
	stateFile    = "integration_test_state.json"
	adminPATFile = "admin-pat.txt"
)

// UserType provides constants that give
// a short explanation with the purpose
// a service user.
// This allows to pre-create users with
// different permissions and reuse them.
type UserType int

//go:generate enumer -type UserType -transform snake -trimprefix UserType
const (
	UserTypeUnspecified UserType = iota
	UserTypeSystem               // UserTypeSystem is a user with access to the system service.
	UserTypeIAMOwner
	UserTypeOrgOwner
	UserTypeLogin
)

const (
	UserPassword = "VeryS3cret!"
)

const (
	PortMilestoneServer = "8081"
	PortQuotaServer     = "8082"
)

// User information with a Personal Access Token.
type User struct {
	ID       string
	Username string
	Token    string
}

type UserMap map[UserType]*User

func (m UserMap) Set(typ UserType, user *User) {
	m[typ] = user
}

func (m UserMap) Get(typ UserType) *User {
	return m[typ]
}

type Config struct {
	Hostname     string
	Port         uint16
	Secure       bool
	LoginURLV2   string
	LogoutURLV2  string
	WebAuthNName string
}

// Host returns the primary host of zitadel, on which the first instance is served.
// http://localhost:8080 by default
func (c *Config) Host() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

// Instance is a Zitadel server and client with all resources available for testing.
type Instance struct {
	Config     Config
	Domain     string
	Instance   *instance.InstanceDetail
	DefaultOrg *org.Org
	Users      UserMap

	Client   *Client
	WebAuthN *webauthn.Client
}

// InitFirstInstance parses config, creates machine users and
// gets default instance and org information.
// Needed details are stored in a state file and can be loaded
// with [loadStateFile] for reuse between multiple test packages.
//
// If an existing state file has the same first instance ID as reported
// by the server, the file will not be modified.
//
// The integration test server must be running.
func InitFirstInstance(ctx context.Context) error {
	var config Config
	if err := yaml.Unmarshal(clientYAML, &config); err != nil {
		return err
	}
	i := &Instance{
		Config: config,
		Domain: config.Hostname,
	}
	token, err := loadInstanceOwnerPAT()
	if err != nil {
		return err
	}
	if err := i.setupInstance(ctx, token); err != nil {
		return err
	}
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(tmpDir, stateFile), data, os.ModePerm)
}

// FirstInstance constructs the first Instance from a reusable state file.
// The returned Instance contains a gRPC client connected to the domain of the new instance.
// The included users are the (global) system user, IAM_OWNER, ORG_OWNER of the default org and
// a Login client user.
// The integration test server must be running.
func FirstInstance(ctx context.Context) (*Instance, error) {
	instance, err := loadStateFile()
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(clientYAML, &instance.Config); err != nil {
		return nil, err
	}

	// refresh short-lived system user token
	if err = instance.createSystemUser(); err != nil {
		return nil, err
	}
	instance.createWebAuthNClient()
	instance.Client, err = newClient(ctx, instance.Host())
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// UseIsolatedInstance creates a new ZITADEL instance with machine users, using the system API.
// The returned Instance contains a gRPC client connected to the domain of the new instance.
// The included users are the (global) system user, IAM_OWNER, ORG_OWNER of the default org and
// a Login client user.
//
// Individual Test function that use an Isolated Instance should use [t.Parallel].
func (i *Instance) UseIsolatedInstance(ctx context.Context) (*Instance, error) {
	systemCtx := i.WithAuthorization(ctx, UserTypeSystem)
	primaryDomain := RandString(5) + ".integration.localhost"
	instance, err := i.Client.System.CreateInstance(systemCtx, &system.CreateInstanceRequest{
		InstanceName: "testinstance",
		CustomDomain: primaryDomain,
		Owner: &system.CreateInstanceRequest_Machine_{
			Machine: &system.CreateInstanceRequest_Machine{
				UserName:            "owner",
				Name:                "owner",
				PersonalAccessToken: &system.CreateInstanceRequest_PersonalAccessToken{},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	newI := &Instance{
		Config: i.Config,
		Domain: primaryDomain,
	}
	if err := newI.setupInstance(ctx, instance.GetPat()); err != nil {
		return nil, err
	}
	newI.createWebAuthNClient()
	_, err = newI.Client.Mgmt.ImportHumanUser(newI.WithAuthorization(ctx, UserTypeIAMOwner), &mgmt.ImportHumanUserRequest{
		UserName: "zitadel-admin@zitadel.localhost",
		Email: &mgmt.ImportHumanUserRequest_Email{
			Email:           "zitadel-admin@zitadel.localhost",
			IsEmailVerified: true,
		},
		Password: "Password1!",
		Profile: &mgmt.ImportHumanUserRequest_Profile{
			FirstName: "hodor",
			LastName:  "hodor",
			NickName:  "hodor",
		},
	})
	if err != nil {
		return nil, err
	}
	return newI, nil
}

func (i *Instance) setupInstance(ctx context.Context, token string) error {
	i.Users = make(UserMap)
	ctx = WithAuthorizationToken(ctx, token)
	return errors.Join(
		i.setClient(ctx),
		i.setInstance(ctx),
		i.setOrganization(ctx),
		i.createSystemUser(),
		i.createMachineUserInstanceOwner(ctx, token),
		i.createMachineUserOrgOwner(ctx),
		i.createLoginClient(ctx),
	)
}

// loadStateFile loads a state file with instance, org and machine user details.
func loadStateFile() (*Instance, error) {
	data, err := os.ReadFile(path.Join(tmpDir, stateFile))
	if err != nil {
		return nil, fmt.Errorf("integration load tester: %w", err)
	}
	dst := new(Instance)
	if err = json.Unmarshal(data, dst); err != nil {
		return nil, fmt.Errorf("integration load tester: %w", err)
	}
	return dst, nil
}

type jsonTester struct {
	Domain       string
	Instance     json.RawMessage
	Organization json.RawMessage
	Users        UserMap
}

func (i *Instance) MarshalJSON() ([]byte, error) {
	instance, err := protojson.Marshal(i.Instance)
	if err != nil {
		return nil, err
	}
	org, err := protojson.Marshal(i.DefaultOrg)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonTester{
		Domain:       i.Domain,
		Instance:     instance,
		Organization: org,
		Users:        i.Users,
	})
}

func (i *Instance) UnmarshalJSON(data []byte) error {
	dst := new(jsonTester)
	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}

	instance := new(instance.InstanceDetail)
	if err := protojson.Unmarshal(dst.Instance, instance); err != nil {
		return err
	}
	org := new(org.Org)
	if err := protojson.Unmarshal(dst.Organization, org); err != nil {
		return err
	}
	*i = Instance{
		Domain:     dst.Domain,
		Instance:   instance,
		DefaultOrg: org,
		Users:      dst.Users,
	}
	return nil
}

// Host returns the primary Domain of the instance with the port.
func (i *Instance) Host() string {
	return fmt.Sprintf("%s:%d", i.Domain, i.Config.Port)
}

func (i *Instance) createSystemUser() error {
	const ISSUER = "tester"
	audience := http_util.BuildOrigin(i.Config.Host(), false)
	signer, err := client.NewSignerFromPrivateKeyByte(systemUserKey, "")
	if err != nil {
		return err
	}
	jwt, err := client.SignedJWTProfileAssertion(ISSUER, []string{audience}, time.Hour, signer)
	if err != nil {
		return err
	}
	i.Users.Set(UserTypeSystem, &User{
		ID:       "SYSTEM",
		Username: "SYSTEM",
		Token:    jwt,
	})
	return nil
}

func loadInstanceOwnerPAT() (string, error) {
	data, err := os.ReadFile(filepath.Join(tmpDir, adminPATFile))
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(data)), nil
}

func (i *Instance) createMachineUserInstanceOwner(ctx context.Context, token string) error {
	user, err := i.Client.Auth.GetMyUser(WithAuthorizationToken(ctx, token), &auth.GetMyUserRequest{})
	if err != nil {
		return err
	}
	i.Users.Set(UserTypeIAMOwner, &User{
		ID:       user.GetUser().GetId(),
		Username: user.GetUser().GetUserName(),
		Token:    token,
	})
	return nil
}

func (i *Instance) createMachineUserOrgOwner(ctx context.Context) error {
	userID, err := i.createMachineUser(ctx, UserTypeOrgOwner)
	if err != nil {
		return err
	}
	_, err = i.Client.Mgmt.AddOrgMember(ctx, &management.AddOrgMemberRequest{
		UserId: userID,
		Roles:  []string{"ORG_OWNER"},
	})
	return err
}

func (i *Instance) createLoginClient(ctx context.Context) error {
	_, err := i.createMachineUser(ctx, UserTypeLogin)
	return err
}

func (i *Instance) setClient(ctx context.Context) error {
	client, err := newClient(ctx, i.Host())
	if err != nil {
		return err
	}
	i.Client = client
	return nil
}

func (i *Instance) setInstance(ctx context.Context) error {
	instance, err := i.Client.Admin.GetMyInstance(ctx, &admin.GetMyInstanceRequest{})
	if err != nil {
		return err
	}
	i.Instance = instance.GetInstance()
	return nil
}

func (i *Instance) setOrganization(ctx context.Context) error {
	resp, err := i.Client.Mgmt.GetMyOrg(ctx, &management.GetMyOrgRequest{})
	if err != nil {
		return err
	}
	i.DefaultOrg = resp.GetOrg()
	return nil
}

func (i *Instance) createMachineUser(ctx context.Context, userType UserType) (userID string, err error) {
	username := gofakeit.Username()
	userResp, err := i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		UserName:        username,
		Name:            username,
		Description:     userType.String(),
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return "", err
	}
	userID = userResp.GetUserId()
	patResp, err := i.Client.Mgmt.AddPersonalAccessToken(ctx, &management.AddPersonalAccessTokenRequest{
		UserId: userID,
	})
	if err != nil {
		return "", err
	}
	i.Users.Set(userType, &User{
		ID:       userID,
		Username: username,
		Token:    patResp.GetToken(),
	})
	return userID, nil
}

func (i *Instance) createWebAuthNClient() {
	i.WebAuthN = webauthn.NewClient(i.Config.WebAuthNName, i.Domain, http_util.BuildOrigin(i.Host(), i.Config.Secure))
}

func (i *Instance) WithAuthorization(ctx context.Context, u UserType) context.Context {
	return i.WithInstanceAuthorization(ctx, u)
}

func (i *Instance) WithInstanceAuthorization(ctx context.Context, u UserType) context.Context {
	return WithAuthorizationToken(ctx, i.Users.Get(u).Token)
}

func (i *Instance) GetUserID(u UserType) string {
	return i.Users.Get(u).ID
}

func WithAuthorizationToken(ctx context.Context, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}
	md.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return metadata.NewOutgoingContext(ctx, md)
}

func (i *Instance) BearerToken(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	return md.Get("Authorization")[0]
}

func (i *Instance) WithSystemAuthorizationHTTP(u UserType) map[string]string {
	return map[string]string{"Authorization": fmt.Sprintf("Bearer %s", i.Users.Get(u).Token)}
}

func runMilestoneServer(ctx context.Context, bodies chan []byte) (*httptest.Server, error) {
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.Header.Get("single-value") != "single-value" {
			http.Error(w, "single-value header not set", http.StatusInternalServerError)
			return
		}
		if reflect.DeepEqual(r.Header.Get("multi-value"), "multi-value-1,multi-value-2") {
			http.Error(w, "single-value header not set", http.StatusInternalServerError)
			return
		}
		bodies <- body
		w.WriteHeader(http.StatusOK)
	}))
	config := net.ListenConfig()
	listener, err := config.Listen(ctx, "tcp", ":"+PortMilestoneServer)
	if err != nil {
		return nil, err
	}
	mockServer.Listener = listener
	mockServer.Start()
	return mockServer, nil
}

func runQuotaServer(ctx context.Context, bodies chan []byte) (*httptest.Server, error) {
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bodies <- body
		w.WriteHeader(http.StatusOK)
	}))
	config := net.ListenConfig()
	listener, err := config.Listen(ctx, "tcp", ":"+PortQuotaServer)
	if err != nil {
		return nil, err
	}
	mockServer.Listener = listener
	mockServer.Start()
	return mockServer, nil
}
