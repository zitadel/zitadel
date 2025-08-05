package query

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"slices"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Users struct {
	SearchResponse
	Users []*User
}

type User struct {
	ID                 string                     `json:"id,omitempty"`
	CreationDate       time.Time                  `json:"creation_date,omitempty"`
	ChangeDate         time.Time                  `json:"change_date,omitempty"`
	ResourceOwner      string                     `json:"resource_owner,omitempty"`
	Sequence           uint64                     `json:"sequence,omitempty"`
	State              domain.UserState           `json:"state,omitempty"`
	Type               domain.UserType            `json:"type,omitempty"`
	Username           string                     `json:"username,omitempty"`
	LoginNames         database.TextArray[string] `json:"login_names,omitempty"`
	PreferredLoginName string                     `json:"preferred_login_name,omitempty"`
	Human              *Human                     `json:"human,omitempty"`
	Machine            *Machine                   `json:"machine,omitempty"`
}

type Human struct {
	FirstName              string              `json:"first_name,omitempty"`
	LastName               string              `json:"last_name,omitempty"`
	NickName               string              `json:"nick_name,omitempty"`
	DisplayName            string              `json:"display_name,omitempty"`
	AvatarKey              string              `json:"avatar_key,omitempty"`
	PreferredLanguage      language.Tag        `json:"preferred_language,omitempty"`
	Gender                 domain.Gender       `json:"gender,omitempty"`
	Email                  domain.EmailAddress `json:"email,omitempty"`
	IsEmailVerified        bool                `json:"is_email_verified,omitempty"`
	Phone                  domain.PhoneNumber  `json:"phone,omitempty"`
	IsPhoneVerified        bool                `json:"is_phone_verified,omitempty"`
	PasswordChangeRequired bool                `json:"password_change_required,omitempty"`
	PasswordChanged        time.Time           `json:"password_changed,omitempty"`
	MFAInitSkipped         time.Time           `json:"mfa_init_skipped,omitempty"`
}

type Profile struct {
	ID                string
	CreationDate      time.Time
	ChangeDate        time.Time
	ResourceOwner     string
	Sequence          uint64
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	AvatarKey         string
	PreferredLanguage language.Tag
	Gender            domain.Gender
}

type Email struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Email         domain.EmailAddress
	IsVerified    bool
}

type Phone struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Phone         string
	IsVerified    bool
}

type Machine struct {
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	EncodedSecret   string               `json:"encoded_hash,omitempty"`
	AccessTokenType domain.OIDCTokenType `json:"access_token_type,omitempty"`
}

type NotifyUser struct {
	ID                 string
	CreationDate       time.Time
	ChangeDate         time.Time
	ResourceOwner      string
	Sequence           uint64
	State              domain.UserState
	Type               domain.UserType
	Username           string
	LoginNames         database.TextArray[string]
	PreferredLoginName string
	FirstName          string
	LastName           string
	NickName           string
	DisplayName        string
	AvatarKey          string
	PreferredLanguage  language.Tag
	Gender             domain.Gender
	LastEmail          string
	VerifiedEmail      string
	LastPhone          string
	VerifiedPhone      string
	PasswordSet        bool
}

func usersCheckPermission(ctx context.Context, users *Users, permissionCheck domain.PermissionCheck) {
	users.Users = slices.DeleteFunc(users.Users,
		func(user *User) bool {
			return userCheckPermission(ctx, user.ResourceOwner, user.ID, permissionCheck) != nil
		},
	)
}

func userPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, filters []SearchQuery) sq.SelectBuilder {
	return userPermissionCheckV2WithCustomColumns(ctx, query, enabled, filters, UserResourceOwnerCol, UserIDCol)
}

func userPermissionCheckV2WithCustomColumns(ctx context.Context, query sq.SelectBuilder, enabled bool, filters []SearchQuery, userResourceOwnerCol, userID Column) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		userResourceOwnerCol,
		domain.PermissionUserRead,
		SingleOrgPermissionOption(filters),
		OwnedRowsPermissionOption(userID),
	)
	return query.JoinClause(join, args...)
}

type UserSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	userTable = table{
		name:          projection.UserTable,
		instanceIDCol: projection.UserInstanceIDCol,
	}
	UserIDCol = Column{
		name:  projection.UserIDCol,
		table: userTable,
	}
	UserCreationDateCol = Column{
		name:  projection.UserCreationDateCol,
		table: userTable,
	}
	UserChangeDateCol = Column{
		name:  projection.UserChangeDateCol,
		table: userTable,
	}
	UserResourceOwnerCol = Column{
		name:  projection.UserResourceOwnerCol,
		table: userTable,
	}
	UserInstanceIDCol = Column{
		name:  projection.UserInstanceIDCol,
		table: userTable,
	}
	UserStateCol = Column{
		name:  projection.UserStateCol,
		table: userTable,
	}
	UserSequenceCol = Column{
		name:  projection.UserSequenceCol,
		table: userTable,
	}
	UserUsernameCol = Column{
		name:           projection.UserUsernameCol,
		table:          userTable,
		isOrderByLower: true,
	}
	UserTypeCol = Column{
		name:  projection.UserTypeCol,
		table: userTable,
	}

	userLoginNamesTable         = loginNameTable.setAlias("login_names")
	userLoginNamesUserIDCol     = LoginNameUserIDCol.setTable(userLoginNamesTable)
	userLoginNamesInstanceIDCol = LoginNameInstanceIDCol.setTable(userLoginNamesTable)
	userLoginNamesListCol       = Column{
		name:  "login_names",
		table: userLoginNamesTable,
	}
	userPreferredLoginNameCol = Column{
		name:  "preferred_login_name",
		table: userLoginNamesTable,
	}
)

var (
	humanTable = table{
		name:          projection.UserHumanTable,
		instanceIDCol: projection.HumanUserInstanceIDCol,
	}
	// profile
	HumanUserIDCol = Column{
		name:  projection.HumanUserIDCol,
		table: humanTable,
	}
	HumanFirstNameCol = Column{
		name:           projection.HumanFirstNameCol,
		table:          humanTable,
		isOrderByLower: true,
	}
	HumanLastNameCol = Column{
		name:           projection.HumanLastNameCol,
		table:          humanTable,
		isOrderByLower: true,
	}
	HumanNickNameCol = Column{
		name:           projection.HumanNickNameCol,
		table:          humanTable,
		isOrderByLower: true,
	}
	HumanDisplayNameCol = Column{
		name:           projection.HumanDisplayNameCol,
		table:          humanTable,
		isOrderByLower: true,
	}
	HumanPreferredLanguageCol = Column{
		name:  projection.HumanPreferredLanguageCol,
		table: humanTable,
	}
	HumanGenderCol = Column{
		name:  projection.HumanGenderCol,
		table: humanTable,
	}
	HumanAvatarURLCol = Column{
		name:  projection.HumanAvatarURLCol,
		table: humanTable,
	}

	// email
	HumanEmailCol = Column{
		name:           projection.HumanEmailCol,
		table:          humanTable,
		isOrderByLower: true,
	}
	HumanIsEmailVerifiedCol = Column{
		name:  projection.HumanIsEmailVerifiedCol,
		table: humanTable,
	}

	// phone
	HumanPhoneCol = Column{
		name:  projection.HumanPhoneCol,
		table: humanTable,
	}
	HumanIsPhoneVerifiedCol = Column{
		name:  projection.HumanIsPhoneVerifiedCol,
		table: humanTable,
	}

	HumanPasswordChangeRequiredCol = Column{
		name:  projection.HumanPasswordChangeRequired,
		table: humanTable,
	}
	HumanPasswordChangedCol = Column{
		name:  projection.HumanPasswordChanged,
		table: humanTable,
	}
	HumanMFAInitSkippedCol = Column{
		name:  projection.HumanMFAInitSkipped,
		table: humanTable,
	}
)

var (
	machineTable = table{
		name:          projection.UserMachineTable,
		instanceIDCol: projection.MachineUserInstanceIDCol,
	}
	MachineUserIDCol = Column{
		name:  projection.MachineUserIDCol,
		table: machineTable,
	}
	MachineNameCol = Column{
		name:           projection.MachineNameCol,
		table:          machineTable,
		isOrderByLower: true,
	}
	MachineDescriptionCol = Column{
		name:  projection.MachineDescriptionCol,
		table: machineTable,
	}
	MachineSecretCol = Column{
		name:  projection.MachineSecretCol,
		table: machineTable,
	}
	MachineAccessTokenTypeCol = Column{
		name:  projection.MachineAccessTokenTypeCol,
		table: machineTable,
	}
)

var (
	notifyTable = table{
		name:          projection.UserNotifyTable,
		instanceIDCol: projection.NotifyInstanceIDCol,
	}
	NotifyUserIDCol = Column{
		name:  projection.NotifyUserIDCol,
		table: notifyTable,
	}
	NotifyEmailCol = Column{
		name:           projection.NotifyLastEmailCol,
		table:          notifyTable,
		isOrderByLower: true,
	}
	NotifyVerifiedEmailCol = Column{
		name:           projection.NotifyVerifiedEmailCol,
		table:          notifyTable,
		isOrderByLower: true,
	}
	NotifyVerifiedEmailLowerCaseCol = Column{
		name:  projection.NotifyVerifiedEmailLowerCol,
		table: notifyTable,
	}
	NotifyPhoneCol = Column{
		name:  projection.NotifyLastPhoneCol,
		table: notifyTable,
	}
	NotifyVerifiedPhoneCol = Column{
		name:  projection.NotifyVerifiedPhoneCol,
		table: notifyTable,
	}
	NotifyPasswordSetCol = Column{
		name:  projection.NotifyPasswordSetCol,
		table: notifyTable,
	}
)

//go:embed user_by_id.sql
var userByIDQuery string

func userCheckPermission(ctx context.Context, resourceOwner string, userID string, permissionCheck domain.PermissionCheck) error {
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := permissionCheck(ctx, domain.PermissionUserRead, resourceOwner, userID); err != nil {
			return err
		}
	}
	return nil
}

func (q *Queries) GetUserByIDWithPermission(ctx context.Context, shouldTriggerBulk bool, userID string, permissionCheck domain.PermissionCheck) (*User, error) {
	user, err := q.GetUserByID(ctx, shouldTriggerBulk, userID)
	if err != nil {
		return nil, err
	}
	if err := userCheckPermission(ctx, user.ResourceOwner, user.ID, permissionCheck); err != nil {
		return nil, err
	}
	return user, nil
}

func (q *Queries) GetUserByID(ctx context.Context, shouldTriggerBulk bool, userID string) (user *User, err error) {
	return q.GetUserByIDWithResourceOwner(ctx, shouldTriggerBulk, userID, "")
}

func (q *Queries) GetUserByIDWithResourceOwner(ctx context.Context, shouldTriggerBulk bool, userID, resourceOwner string) (user *User, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		triggerUserProjections(ctx)
	}

	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			user, err = scanUser(row)
			return err
		},
		userByIDQuery,
		userID,
		resourceOwner,
		authz.GetInstance(ctx).InstanceID(),
	)
	return user, err
}

//go:embed user_by_login_name.sql
var userByLoginNameQuery string

func (q *Queries) GetUserByLoginName(ctx context.Context, shouldTriggered bool, loginName string) (user *User, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggered {
		triggerUserProjections(ctx)
	}

	loginName = strings.ToLower(loginName)

	username := loginName
	domainIndex := strings.LastIndex(loginName, "@")
	var domainSuffix string
	// split between the last @ (so ignore it if the login name ends with it)
	if domainIndex > 0 && domainIndex != len(loginName)-1 {
		domainSuffix = loginName[domainIndex+1:]
		username = loginName[:domainIndex]
	}

	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			user, err = scanUser(row)
			return err
		},
		userByLoginNameQuery,
		username,
		domainSuffix,
		loginName,
		authz.GetInstance(ctx).InstanceID(),
	)
	return user, err
}

func (q *Queries) GetHumanProfile(ctx context.Context, userID string, queries ...SearchQuery) (profile *Profile, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProfileQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		UserIDCol.identifier():         userID,
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgbg2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		profile, err = scan(row)
		return err
	}, stmt, args...)
	return profile, err
}

func (q *Queries) GetHumanEmail(ctx context.Context, userID string, queries ...SearchQuery) (email *Email, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareEmailQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		UserIDCol.identifier():         userID,
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-BHhj3", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		email, err = scan(row)
		return err
	}, stmt, args...)
	return email, err
}

func (q *Queries) GetHumanPhone(ctx context.Context, userID string, queries ...SearchQuery) (phone *Phone, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := preparePhoneQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		UserIDCol.identifier():         userID,
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dg43g", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		phone, err = scan(row)
		return err
	}, stmt, args...)
	return phone, err
}

//go:embed user_notify_by_id.sql
var notifyUserByIDQuery string

func (q *Queries) GetNotifyUserByID(ctx context.Context, shouldTriggered bool, userID string) (user *NotifyUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggered {
		triggerUserProjections(ctx)
	}

	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			user, err = scanNotifyUser(row)
			return err
		},
		notifyUserByIDQuery,
		userID,
		authz.GetInstance(ctx).InstanceID(),
	)
	return user, err
}

//go:embed user_notify_by_login_name.sql
var notifyUserByLoginNameQuery string

func (q *Queries) GetNotifyUserByLoginName(ctx context.Context, shouldTriggered bool, loginName string) (user *NotifyUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggered {
		triggerUserProjections(ctx)
	}

	loginName = strings.ToLower(loginName)

	username := loginName
	domainIndex := strings.LastIndex(loginName, "@")
	var domainSuffix string
	// split between the last @ (so ignore it if the login name ends with it)
	if domainIndex > 0 && domainIndex != len(loginName)-1 {
		domainSuffix = loginName[domainIndex+1:]
		username = loginName[:domainIndex]
	}

	err = q.client.QueryRowContext(ctx,
		func(row *sql.Row) error {
			user, err = scanNotifyUser(row)
			return err
		},
		notifyUserByLoginNameQuery,
		username,
		domainSuffix,
		loginName,
		authz.GetInstance(ctx).InstanceID(),
	)
	return user, err
}

func (q *Queries) GetNotifyUser(ctx context.Context, shouldTriggered bool, queries ...SearchQuery) (user *NotifyUser, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggered {
		triggerUserProjections(ctx)
	}

	query, scan := prepareNotifyUserQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Err3g", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		user, err = scan(row)
		return err
	}, stmt, args...)
	return user, err
}

func (q *Queries) CountUsers(ctx context.Context, queries *UserSearchQueries) (count uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareCountUsersQuery()
	eq := sq.Eq{UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return 0, zerrors.ThrowInternal(err, "QUERY-w3Dx", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		count, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return 0, zerrors.ThrowInternal(err, "QUERY-AG4gs", "Errors.Internal")
	}
	return count, err
}

func (q *Queries) SearchUsers(ctx context.Context, queries *UserSearchQueries, permissionCheck domain.PermissionCheck) (*Users, error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	users, err := q.searchUsers(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		usersCheckPermission(ctx, users, permissionCheck)
	}
	return users, nil
}

func (q *Queries) searchUsers(ctx context.Context, queries *UserSearchQueries, permissionCheckV2 bool) (users *Users, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUsersQuery()
	query = userPermissionCheckV2(ctx, query, permissionCheckV2, queries.Queries)
	stmt, args, err := queries.toQuery(query).Where(sq.Eq{
		UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Dgbg2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		users, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-AG4gs", "Errors.Internal")
	}
	users.State, err = q.latestState(ctx, userTable)
	return users, err
}

func (q *Queries) IsUserUnique(ctx context.Context, username, email, resourceOwner string) (isUnique bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareUserUniqueQuery()
	queries := make([]SearchQuery, 0, 3)
	if username != "" {
		usernameQuery, err := NewUserUsernameSearchQuery(username, TextEquals)
		if err != nil {
			return false, err
		}
		queries = append(queries, usernameQuery)
	}
	if email != "" {
		emailQuery, err := NewUserEmailSearchQuery(email, TextEquals)
		if err != nil {
			return false, err
		}
		queries = append(queries, emailQuery)
	}
	if resourceOwner != "" {
		resourceOwnerQuery, err := NewUserResourceOwnerSearchQuery(resourceOwner, TextEquals)
		if err != nil {
			return false, err
		}
		queries = append(queries, resourceOwnerQuery)
	}
	for _, q := range queries {
		query = q.toQuery(query)
	}
	eq := sq.Eq{UserInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return false, zerrors.ThrowInternal(err, "QUERY-Dg43g", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		isUnique, err = scan(row)
		return err
	}, stmt, args...)
	return isUnique, err
}

//go:embed user_claimed_user_ids.sql
var userClaimedUserIDOfOrgDomain string

func (q *Queries) SearchClaimedUserIDsOfOrgDomain(ctx context.Context, domain, orgID string) (userIDs []string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			userIDs = make([]string, 0)
			for rows.Next() {
				var userID string
				err := rows.Scan(&userID)
				if err != nil {
					return err
				}
				userIDs = append(userIDs, userID)
			}
			return nil
		},
		userClaimedUserIDOfOrgDomain,
		authz.GetInstance(ctx).InstanceID(),
		"%@"+domain,
		orgID,
	)

	return userIDs, err
}

func (q *UserSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (r *UserSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewUserResourceOwnerSearchQuery(orgID, TextEquals)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func NewUserOrSearchQuery(values []SearchQuery) (SearchQuery, error) {
	return NewOrQuery(values...)
}

func NewUserAndSearchQuery(values []SearchQuery) (SearchQuery, error) {
	return NewAndQuery(values...)
}

func NewUserNotSearchQuery(value SearchQuery) (SearchQuery, error) {
	return NewNotQuery(value)
}

func NewUserInUserIdsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(UserIDCol, values)
}

func NewUserInUserEmailsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(HumanEmailCol, values)
}

func NewUserResourceOwnerSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserResourceOwnerCol, value, comparison)
}

func NewUserUsernameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(UserUsernameCol, value, comparison)
}

func NewUserFirstNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanFirstNameCol, value, comparison)
}

func NewUserLastNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanLastNameCol, value, comparison)
}

func NewUserNickNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanNickNameCol, value, comparison)
}

func NewUserDisplayNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanDisplayNameCol, value, comparison)
}

func NewUserEmailSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanEmailCol, value, comparison)
}

func NewUserPhoneSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(HumanPhoneCol, value, comparison)
}

func NewUserVerifiedEmailSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(NotifyVerifiedEmailLowerCaseCol, strings.ToLower(value), TextEquals)
}

func NewUserVerifiedPhoneSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(NotifyVerifiedPhoneCol, value, comparison)
}

func NewUserStateSearchQuery(value domain.UserState) (SearchQuery, error) {
	return NewNumberQuery(UserStateCol, value, NumberEquals)
}

func NewUserTypeSearchQuery(value domain.UserType) (SearchQuery, error) {
	return NewNumberQuery(UserTypeCol, value, NumberEquals)
}

func NewUserPreferredLoginNameSearchQuery(value string, comparison TextComparison) (SearchQuery, error) {
	return NewTextQuery(userPreferredLoginNameCol, value, comparison)
}

func NewUserLoginNameExistsQuery(value string, comparison TextComparison) (SearchQuery, error) {
	// linking queries for the sub select
	instanceQuery, err := NewColumnComparisonQuery(LoginNameInstanceIDCol, UserInstanceIDCol, ColumnEquals)
	if err != nil {
		return nil, err
	}
	userIDQuery, err := NewColumnComparisonQuery(LoginNameUserIDCol, UserIDCol, ColumnEquals)
	if err != nil {
		return nil, err
	}
	// text query to select data from the linked sub select
	loginNameQuery, err := NewTextQuery(LoginNameNameCol, value, comparison)
	if err != nil {
		return nil, err
	}
	// full definition of the sub select
	subSelect, err := NewSubSelect(LoginNameUserIDCol, []SearchQuery{instanceQuery, userIDQuery, loginNameQuery})
	if err != nil {
		return nil, err
	}
	// "WHERE * IN (*)" query with subquery as list-data provider
	return NewListQuery(
		UserIDCol,
		subSelect,
		ListIn,
	)
}

func triggerUserProjections(ctx context.Context) {
	triggerBatch(ctx, projection.UserProjection, projection.LoginNameProjection)
}

var joinLoginNames = `LEFT JOIN LATERAL (` +
	`SELECT` +
	` ARRAY_AGG(ln.login_name ORDER BY ln.login_name) AS login_names,` +
	` MAX(CASE WHEN ln.is_primary THEN ln.login_name ELSE NULL END) AS preferred_login_name` +
	` FROM` +
	` projections.login_names3 AS ln` +
	` WHERE` +
	` ln.user_id = ` + UserIDCol.identifier() +
	` AND ln.instance_id = ` + UserInstanceIDCol.identifier() +
	`) AS login_names ON TRUE`

func scanUser(row *sql.Row) (*User, error) {
	u := new(User)
	var count int
	preferredLoginName := sql.NullString{}

	humanID := sql.NullString{}
	firstName := sql.NullString{}
	lastName := sql.NullString{}
	nickName := sql.NullString{}
	displayName := sql.NullString{}
	preferredLanguage := sql.NullString{}
	gender := sql.NullInt32{}
	avatarKey := sql.NullString{}
	email := sql.NullString{}
	isEmailVerified := sql.NullBool{}
	phone := sql.NullString{}
	isPhoneVerified := sql.NullBool{}
	passwordChangeRequired := sql.NullBool{}
	passwordChanged := sql.NullTime{}
	mfaInitSkipped := sql.NullTime{}

	machineID := sql.NullString{}
	name := sql.NullString{}
	description := sql.NullString{}
	encodedHash := sql.NullString{}
	accessTokenType := sql.NullInt32{}

	err := row.Scan(
		&u.ID,
		&u.CreationDate,
		&u.ChangeDate,
		&u.ResourceOwner,
		&u.Sequence,
		&u.State,
		&u.Type,
		&u.Username,
		&u.LoginNames,
		&preferredLoginName,
		&humanID,
		&firstName,
		&lastName,
		&nickName,
		&displayName,
		&preferredLanguage,
		&gender,
		&avatarKey,
		&email,
		&isEmailVerified,
		&phone,
		&isPhoneVerified,
		&passwordChangeRequired,
		&passwordChanged,
		&mfaInitSkipped,
		&machineID,
		&name,
		&description,
		&encodedHash,
		&accessTokenType,
		&count,
	)

	if err != nil || count != 1 {
		if errors.Is(err, sql.ErrNoRows) || count != 1 {
			return nil, zerrors.ThrowNotFound(err, "QUERY-Dfbg2", "Errors.User.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-Bgah2", "Errors.Internal")
	}

	u.PreferredLoginName = preferredLoginName.String

	if humanID.Valid {
		u.Human = &Human{
			FirstName:              firstName.String,
			LastName:               lastName.String,
			NickName:               nickName.String,
			DisplayName:            displayName.String,
			AvatarKey:              avatarKey.String,
			PreferredLanguage:      language.Make(preferredLanguage.String),
			Gender:                 domain.Gender(gender.Int32),
			Email:                  domain.EmailAddress(email.String),
			IsEmailVerified:        isEmailVerified.Bool,
			Phone:                  domain.PhoneNumber(phone.String),
			IsPhoneVerified:        isPhoneVerified.Bool,
			PasswordChangeRequired: passwordChangeRequired.Bool,
			PasswordChanged:        passwordChanged.Time,
			MFAInitSkipped:         mfaInitSkipped.Time,
		}
	} else if machineID.Valid {
		u.Machine = &Machine{
			Name:            name.String,
			Description:     description.String,
			EncodedSecret:   encodedHash.String,
			AccessTokenType: domain.OIDCTokenType(accessTokenType.Int32),
		}
	}
	return u, nil
}

func prepareProfileQuery() (sq.SelectBuilder, func(*sql.Row) (*Profile, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanNickNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanPreferredLanguageCol.identifier(),
			HumanGenderCol.identifier(),
			HumanAvatarURLCol.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Profile, error) {
			p := new(Profile)

			humanID := sql.NullString{}
			firstName := sql.NullString{}
			lastName := sql.NullString{}
			nickName := sql.NullString{}
			displayName := sql.NullString{}
			preferredLanguage := sql.NullString{}
			gender := sql.NullInt32{}
			avatarKey := sql.NullString{}
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
				&humanID,
				&firstName,
				&lastName,
				&nickName,
				&displayName,
				&preferredLanguage,
				&gender,
				&avatarKey,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-HNhb3", "Errors.User.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Rfheq", "Errors.Internal")
			}
			if !humanID.Valid {
				return nil, zerrors.ThrowPreconditionFailed(nil, "QUERY-WLTce", "Errors.User.NotHuman")
			}

			p.FirstName = firstName.String
			p.LastName = lastName.String
			p.NickName = nickName.String
			p.DisplayName = displayName.String
			p.AvatarKey = avatarKey.String
			p.PreferredLanguage = language.Make(preferredLanguage.String)
			p.Gender = domain.Gender(gender.Int32)

			return p, nil
		}
}

func prepareEmailQuery() (sq.SelectBuilder, func(*sql.Row) (*Email, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanEmailCol.identifier(),
			HumanIsEmailVerifiedCol.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Email, error) {
			e := new(Email)

			humanID := sql.NullString{}
			email := sql.NullString{}
			isEmailVerified := sql.NullBool{}

			err := row.Scan(
				&e.ID,
				&e.CreationDate,
				&e.ChangeDate,
				&e.ResourceOwner,
				&e.Sequence,
				&humanID,
				&email,
				&isEmailVerified,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-Hms2s", "Errors.User.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Nu42d", "Errors.Internal")
			}
			if !humanID.Valid {
				return nil, zerrors.ThrowPreconditionFailed(nil, "QUERY-pt7HY", "Errors.User.NotHuman")
			}

			e.Email = domain.EmailAddress(email.String)
			e.IsVerified = isEmailVerified.Bool

			return e, nil
		}
}

func preparePhoneQuery() (sq.SelectBuilder, func(*sql.Row) (*Phone, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanPhoneCol.identifier(),
			HumanIsPhoneVerifiedCol.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Phone, error) {
			e := new(Phone)

			humanID := sql.NullString{}
			phone := sql.NullString{}
			isPhoneVerified := sql.NullBool{}

			err := row.Scan(
				&e.ID,
				&e.CreationDate,
				&e.ChangeDate,
				&e.ResourceOwner,
				&e.Sequence,
				&humanID,
				&phone,
				&isPhoneVerified,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-DAvb3", "Errors.User.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-Bmf2h", "Errors.Internal")
			}
			if !humanID.Valid {
				return nil, zerrors.ThrowPreconditionFailed(nil, "QUERY-hliQl", "Errors.User.NotHuman")
			}

			e.Phone = phone.String
			e.IsVerified = isPhoneVerified.Bool

			return e, nil
		}
}

func prepareNotifyUserQuery() (sq.SelectBuilder, func(*sql.Row) (*NotifyUser, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			UserStateCol.identifier(),
			UserTypeCol.identifier(),
			UserUsernameCol.identifier(),
			userLoginNamesListCol.identifier(),
			userPreferredLoginNameCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanNickNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanPreferredLanguageCol.identifier(),
			HumanGenderCol.identifier(),
			HumanAvatarURLCol.identifier(),
			NotifyUserIDCol.identifier(),
			NotifyEmailCol.identifier(),
			NotifyVerifiedEmailCol.identifier(),
			NotifyPhoneCol.identifier(),
			NotifyVerifiedPhoneCol.identifier(),
			NotifyPasswordSetCol.identifier(),
			countColumn.identifier(),
		).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			LeftJoin(join(NotifyUserIDCol, UserIDCol)).
			JoinClause(joinLoginNames).
			PlaceholderFormat(sq.Dollar),
		scanNotifyUser
}

func scanNotifyUser(row *sql.Row) (*NotifyUser, error) {
	u := new(NotifyUser)
	var count int
	loginNames := database.TextArray[string]{}
	preferredLoginName := sql.NullString{}

	humanID := sql.NullString{}
	firstName := sql.NullString{}
	lastName := sql.NullString{}
	nickName := sql.NullString{}
	displayName := sql.NullString{}
	preferredLanguage := sql.NullString{}
	gender := sql.NullInt32{}
	avatarKey := sql.NullString{}

	notifyUserID := sql.NullString{}
	notifyEmail := sql.NullString{}
	notifyVerifiedEmail := sql.NullString{}
	notifyPhone := sql.NullString{}
	notifyVerifiedPhone := sql.NullString{}
	notifyPasswordSet := sql.NullBool{}

	err := row.Scan(
		&u.ID,
		&u.CreationDate,
		&u.ChangeDate,
		&u.ResourceOwner,
		&u.Sequence,
		&u.State,
		&u.Type,
		&u.Username,
		&loginNames,
		&preferredLoginName,
		&humanID,
		&firstName,
		&lastName,
		&nickName,
		&displayName,
		&preferredLanguage,
		&gender,
		&avatarKey,
		&notifyUserID,
		&notifyEmail,
		&notifyVerifiedEmail,
		&notifyPhone,
		&notifyVerifiedPhone,
		&notifyPasswordSet,
		&count,
	)

	if err != nil || count != 1 {
		if errors.Is(err, sql.ErrNoRows) || count != 1 {
			return nil, zerrors.ThrowNotFound(err, "QUERY-Dgqd2", "Errors.User.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-Dbwsg", "Errors.Internal")
	}

	if !notifyUserID.Valid {
		return nil, zerrors.ThrowPreconditionFailed(nil, "QUERY-Sfw3f", "Errors.User.NotFound")
	}

	u.LoginNames = loginNames
	if preferredLoginName.Valid {
		u.PreferredLoginName = preferredLoginName.String
	}
	if humanID.Valid {
		u.FirstName = firstName.String
		u.LastName = lastName.String
		u.NickName = nickName.String
		u.DisplayName = displayName.String
		u.AvatarKey = avatarKey.String
		u.PreferredLanguage = language.Make(preferredLanguage.String)
		u.Gender = domain.Gender(gender.Int32)
	}
	u.LastEmail = notifyEmail.String
	u.VerifiedEmail = notifyVerifiedEmail.String
	u.LastPhone = notifyPhone.String
	u.VerifiedPhone = notifyVerifiedPhone.String
	u.PasswordSet = notifyPasswordSet.Bool

	return u, nil
}

func prepareCountUsersQuery() (sq.SelectBuilder, func(*sql.Rows) (uint64, error)) {
	return sq.Select(countColumn.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			LeftJoin(join(MachineUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (count uint64, err error) {
			// the count is implemented as a windowing function,
			// if it is zero, no row is returned at all.
			if !rows.Next() {
				return
			}

			err = rows.Scan(&count)
			return
		}
}

func prepareUserUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserStateCol.identifier(),
			UserUsernameCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanEmailCol.identifier(),
			HumanIsEmailVerifiedCol.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (bool, error) {
			userID := sql.NullString{}
			state := sql.NullInt32{}
			username := sql.NullString{}
			humanID := sql.NullString{}
			email := sql.NullString{}
			isEmailVerified := sql.NullBool{}

			err := row.Scan(
				&userID,
				&state,
				&username,
				&humanID,
				&email,
				&isEmailVerified,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return true, nil
				}
				return false, zerrors.ThrowInternal(err, "QUERY-Cxces", "Errors.Internal")
			}
			return !userID.Valid, nil
		}
}

func prepareUsersQuery() (sq.SelectBuilder, func(*sql.Rows) (*Users, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			UserStateCol.identifier(),
			UserTypeCol.identifier(),
			UserUsernameCol.identifier(),
			userLoginNamesListCol.identifier(),
			userPreferredLoginNameCol.identifier(),
			HumanUserIDCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanNickNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanPreferredLanguageCol.identifier(),
			HumanGenderCol.identifier(),
			HumanAvatarURLCol.identifier(),
			HumanEmailCol.identifier(),
			HumanIsEmailVerifiedCol.identifier(),
			HumanPhoneCol.identifier(),
			HumanIsPhoneVerifiedCol.identifier(),
			HumanPasswordChangeRequiredCol.identifier(),
			HumanPasswordChangedCol.identifier(),
			MachineUserIDCol.identifier(),
			MachineNameCol.identifier(),
			MachineDescriptionCol.identifier(),
			MachineSecretCol.identifier(),
			MachineAccessTokenTypeCol.identifier(),
			countColumn.identifier()).
			From(userTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			LeftJoin(join(MachineUserIDCol, UserIDCol)).
			JoinClause(joinLoginNames).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Users, error) {
			users := make([]*User, 0)
			var count uint64
			for rows.Next() {
				u := new(User)
				loginNames := database.TextArray[string]{}
				preferredLoginName := sql.NullString{}

				humanID := sql.NullString{}
				firstName := sql.NullString{}
				lastName := sql.NullString{}
				nickName := sql.NullString{}
				displayName := sql.NullString{}
				preferredLanguage := sql.NullString{}
				gender := sql.NullInt32{}
				avatarKey := sql.NullString{}
				email := sql.NullString{}
				isEmailVerified := sql.NullBool{}
				phone := sql.NullString{}
				isPhoneVerified := sql.NullBool{}
				passwordChangeRequired := sql.NullBool{}
				passwordChanged := sql.NullTime{}

				machineID := sql.NullString{}
				name := sql.NullString{}
				description := sql.NullString{}
				encodedHash := sql.NullString{}
				accessTokenType := sql.NullInt32{}

				err := rows.Scan(
					&u.ID,
					&u.CreationDate,
					&u.ChangeDate,
					&u.ResourceOwner,
					&u.Sequence,
					&u.State,
					&u.Type,
					&u.Username,
					&loginNames,
					&preferredLoginName,
					&humanID,
					&firstName,
					&lastName,
					&nickName,
					&displayName,
					&preferredLanguage,
					&gender,
					&avatarKey,
					&email,
					&isEmailVerified,
					&phone,
					&isPhoneVerified,
					&passwordChangeRequired,
					&passwordChanged,
					&machineID,
					&name,
					&description,
					&encodedHash,
					&accessTokenType,
					&count,
				)
				if err != nil {
					return nil, err
				}

				u.LoginNames = loginNames
				if preferredLoginName.Valid {
					u.PreferredLoginName = preferredLoginName.String
				}

				if humanID.Valid {
					u.Human = &Human{
						FirstName:              firstName.String,
						LastName:               lastName.String,
						NickName:               nickName.String,
						DisplayName:            displayName.String,
						AvatarKey:              avatarKey.String,
						PreferredLanguage:      language.Make(preferredLanguage.String),
						Gender:                 domain.Gender(gender.Int32),
						Email:                  domain.EmailAddress(email.String),
						IsEmailVerified:        isEmailVerified.Bool,
						Phone:                  domain.PhoneNumber(phone.String),
						IsPhoneVerified:        isPhoneVerified.Bool,
						PasswordChangeRequired: passwordChangeRequired.Bool,
						PasswordChanged:        passwordChanged.Time,
					}
				} else if machineID.Valid {
					u.Machine = &Machine{
						Name:            name.String,
						Description:     description.String,
						EncodedSecret:   encodedHash.String,
						AccessTokenType: domain.OIDCTokenType(accessTokenType.Int32),
					}
				}

				users = append(users, u)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-frhbd", "Errors.Query.CloseRows")
			}

			return &Users{
				Users: users,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
