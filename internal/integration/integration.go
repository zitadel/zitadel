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
	"github.com/zitadel/zitadel/pkg/grpc/org"
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
	FirstInstanceUsersKey = "first"
	UserPassword          = "VeryS3cret!"
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

type InstanceUserMap map[string]map[UserType]*User

func (m InstanceUserMap) Set(instanceID string, typ UserType, user *User) {
	if m[instanceID] == nil {
		m[instanceID] = make(map[UserType]*User)
	}
	m[instanceID][typ] = user
}

func (m InstanceUserMap) Get(instanceID string, typ UserType) *User {
	if users, ok := m[instanceID]; ok {
		return users[typ]
	}
	return nil
}

type Config struct {
	Hostname     string
	Port         uint16
	Secure       bool
	LoginURLV2   string
	LogoutURLV2  string
	WebAuthNName string
}

// Tester is a Zitadel server and client with all resources available for testing.
type Tester struct {
	Config       Config
	Instance     *instance.InstanceDetail
	Organisation *org.Org
	Users        InstanceUserMap

	Client   *Client
	WebAuthN *webauthn.Client
}

func (c *Config) Host() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

// NewTester constructs a new Tester from a reusable state file,
// and constructs the gRPC clients.
// The integration test server must be running.
func NewTester(ctx context.Context) (*Tester, error) {
	tester, err := loadStateFile()
	if err != nil {
		return nil, err
	}
	// refresh short-lived system user token
	if err = tester.createSystemUser(); err != nil {
		return nil, err
	}
	tester.WebAuthN = webauthn.NewClient(tester.Config.WebAuthNName, tester.Config.Hostname, http_util.BuildOrigin(tester.Host(), tester.Config.Secure))
	tester.Client, err = newClient(ctx, tester.Host())
	if err != nil {
		return nil, err
	}
	return tester, nil
}

// loadStateFile loads a state file with instance, org and machine user details.
func loadStateFile() (*Tester, error) {
	data, err := os.ReadFile(path.Join(tmpDir, stateFile))
	if err != nil {
		return nil, fmt.Errorf("integration load tester: %w", err)
	}
	dst := new(Tester)
	if err = json.Unmarshal(data, dst); err != nil {
		return nil, fmt.Errorf("integration load tester: %w", err)
	}
	return dst, nil
}

type jsonTester struct {
	Config       Config
	Instance     json.RawMessage
	Organization json.RawMessage
	Users        InstanceUserMap
}

func (s *Tester) MarshalJSON() ([]byte, error) {
	instance, err := protojson.Marshal(s.Instance)
	if err != nil {
		return nil, err
	}
	org, err := protojson.Marshal(s.Organisation)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonTester{
		Config:       s.Config,
		Instance:     instance,
		Organization: org,
		Users:        s.Users,
	})
}

func (s *Tester) UnmarshalJSON(data []byte) error {
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
	*s = Tester{
		Config:       dst.Config,
		Instance:     instance,
		Organisation: org,
		Users:        dst.Users,
	}
	return nil
}

func (s *Tester) Host() string {
	return s.Config.Host()
}

func (s *Tester) createSystemUser() error {
	const ISSUER = "tester"
	audience := http_util.BuildOrigin(s.Host(), false)
	signer, err := client.NewSignerFromPrivateKeyByte(systemUserKey, "")
	if err != nil {
		return err
	}
	jwt, err := client.SignedJWTProfileAssertion(ISSUER, []string{audience}, time.Hour, signer)
	if err != nil {
		return err
	}
	s.Users.Set(FirstInstanceUsersKey, UserTypeSystem, &User{
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

func (s *Tester) createMachineUserInstanceOwner(ctx context.Context, token string) error {
	user, err := s.Client.Auth.GetMyUser(WithAuthorizationToken(ctx, token), &auth.GetMyUserRequest{})
	if err != nil {
		return err
	}
	s.Users.Set(FirstInstanceUsersKey, UserTypeIAMOwner, &User{
		ID:       user.GetUser().GetId(),
		Username: user.GetUser().GetUserName(),
		Token:    token,
	})
	return nil
}

func (s *Tester) createMachineUserOrgOwner(ctx context.Context) error {
	userID, err := s.createMachineUser(ctx, UserTypeOrgOwner)
	if err != nil {
		return err
	}
	_, err = s.Client.Mgmt.AddOrgMember(ctx, &management.AddOrgMemberRequest{
		UserId: userID,
		Roles:  []string{"ORG_OWNER"},
	})
	return err
}

func (s *Tester) createLoginClient(ctx context.Context) error {
	_, err := s.createMachineUser(ctx, UserTypeLogin)
	return err
}

func (s *Tester) setOrganization(ctx context.Context) error {
	resp, err := s.Client.Mgmt.GetMyOrg(ctx, &management.GetMyOrgRequest{})
	if err != nil {
		return err
	}
	s.Organisation = resp.GetOrg()
	return nil
}

func (s *Tester) createMachineUser(ctx context.Context, userType UserType) (userID string, err error) {
	username := gofakeit.Username()
	userResp, err := s.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		UserName:        username,
		Name:            username,
		Description:     userType.String(),
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return "", err
	}
	userID = userResp.GetUserId()
	patResp, err := s.Client.Mgmt.AddPersonalAccessToken(ctx, &management.AddPersonalAccessTokenRequest{
		UserId: userID,
	})
	if err != nil {
		return "", err
	}
	s.Users.Set(FirstInstanceUsersKey, userType, &User{
		ID:       userID,
		Username: username,
		Token:    patResp.GetToken(),
	})
	return userID, nil
}

func (s *Tester) WithAuthorization(ctx context.Context, u UserType) context.Context {
	return s.WithInstanceAuthorization(ctx, u, FirstInstanceUsersKey)
}

func (s *Tester) WithInstanceAuthorization(ctx context.Context, u UserType, instanceID string) context.Context {
	return WithAuthorizationToken(ctx, s.Users.Get(instanceID, u).Token)
}

func (s *Tester) GetUserID(u UserType) string {
	return s.Users.Get(FirstInstanceUsersKey, u).ID
}

func WithAuthorizationToken(ctx context.Context, token string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = make(metadata.MD)
	}
	md.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return metadata.NewOutgoingContext(ctx, md)
}

func (s *Tester) BearerToken(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}
	return md.Get("Authorization")[0]
}

func (s *Tester) WithSystemAuthorizationHTTP(u UserType) map[string]string {
	return map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.Users.Get(FirstInstanceUsersKey, u).Token)}
}

// InitTesterState parses config, creates machine users and
// gets default instance and org information.
// Needed details are stored in a state file and can be loaded
// with [loadStateFile] for reuse between multiple test packages.
//
// If an existing state file has the same first instance ID as reported
// by the server, the file will not be modified.
//
// The integration test server must be running.
func InitTesterState(ctx context.Context) error {
	var config Config
	if err := yaml.Unmarshal(clientYAML, &config); err != nil {
		return err
	}
	client, err := newClient(ctx, config.Host())
	if err != nil {
		return err
	}
	token, err := loadInstanceOwnerPAT()
	if err != nil {
		return err
	}

	ctx = WithAuthorizationToken(ctx, token)
	instance, err := client.Admin.GetMyInstance(ctx, &admin.GetMyInstanceRequest{})
	if err != nil {
		return err
	}
	tester, err := loadStateFile()
	if err == nil && tester.Instance.GetId() == instance.GetInstance().GetId() {
		return nil
	}
	tester = &Tester{
		Users: make(InstanceUserMap),
	}
	tester.Instance = instance.GetInstance()
	tester.Client = client
	tester.Config = config

	err = errors.Join(
		tester.setOrganization(ctx),
		tester.createSystemUser(),
		tester.createMachineUserInstanceOwner(ctx, token),
		tester.createMachineUserOrgOwner(ctx),
		tester.createLoginClient(ctx),
	)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(tester, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(tmpDir, stateFile), data, os.ModePerm)
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
