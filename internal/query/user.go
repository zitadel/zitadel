package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/query/projection"
)

type Users struct {
	SearchResponse
	Users []User
}

type User interface {
}

type Human struct {
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	AvatarKey         string
	PreferredLanguage language.Tag
	Gender            domain.Gender
	Email             string
	IsEmailVerified   bool
	Phone             string
	IsPhoneVerified   bool
}

type Machine struct {
	Name        string
	Description string
}

type user struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	State         domain.UserState
	Username      string
}

type UserSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	userTable = table{
		name: projection.UserTable,
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
	UserStateCol = Column{
		name:  projection.UserStateCol,
		table: userTable,
	}
	UserSequenceCol = Column{
		name:  projection.UserSequenceCol,
		table: userTable,
	}
	UserUsernameCol = Column{
		name:  projection.UserUsernameCol,
		table: userTable,
	}
	UserTypeCol = Column{
		name:  projection.UserTypeCol,
		table: userTable,
	}
)

var (
	humanTable = table{
		name: projection.UserHumanTable,
	}
	// profile
	HumanUserIDCol = Column{
		name:  projection.HumanUserIDCol,
		table: humanTable,
	}
	HumanFirstNameCol = Column{
		name:  projection.HumanFirstNameCol,
		table: humanTable,
	}
	HumanLastNameCol = Column{
		name:  projection.HumanLastNameCol,
		table: humanTable,
	}
	HumanNickNameCol = Column{
		name:  projection.HumanNickNameCol,
		table: humanTable,
	}
	HumanDisplayNameCol = Column{
		name:  projection.HumanDisplayNameCol,
		table: humanTable,
	}
	HumanPreferredLanguageCol = Column{
		name:  projection.HumanPreferredLanguageCol,
		table: humanTable,
	}
	HumanGenderCol = Column{
		name:  projection.HumanGenderCol,
		table: humanTable,
	}
	HumanAvaterURLCol = Column{
		name:  projection.HumanAvaterURLCol,
		table: humanTable,
	}

	// email
	HumanEmailCol = Column{
		name:  projection.HumanEmailCol,
		table: humanTable,
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
)

var (
	machineTable = table{
		name: projection.UserMachineTable,
	}
	MachineUserIDCol = Column{
		name:  projection.MachineUserIDCol,
		table: machineTable,
	}
	MachineNameCol = Column{
		name:  projection.MachineNameCol,
		table: machineTable,
	}
	MachineDescriptionCol = Column{
		name:  projection.MachineDescriptionCol,
		table: machineTable,
	}
)

//func (q *Queries) GetUserByID(ctx context.Context, userID string, queries ...SearchQuery) (User, error) {
//	return q.GetUser(ctx, append(queries, sq.Eq{
//		UserIDCol.identifier(): userID,
//	})...)
//}

func (q *Queries) GetUser(ctx context.Context, queries ...SearchQuery) (User, error) {
	query, scan := prepareUserQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgbg2", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) SearchUsers(ctx context.Context, queries ...SearchQuery) (*Users, error) {
	query, scan := prepareUsersQuery()
	for _, q := range queries {
		query = q.toQuery(query)
	}
	stmt, args, err := query.ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Dgbg2", "Errors.Query.SQLStatment")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AG4gs", "Errors.Internal")
	}
	users, err := scan(rows)
	if err != nil {
		return nil, err
	}
	users.LatestSequence, err = q.latestSequence(ctx, projectRolesTable)
	return users, err
}

func (q *Queries) IsUserUnique(ctx context.Context, queries ...SearchQuery) (bool, error) {
	return false, nil
}

func (q *Queries) UserEvents(ctx context.Context, orgID, userID string, sequence uint64) ([]eventstore.Event, error) {
	query := NewUserEventSearchQuery(userID, orgID, sequence)
	return q.eventstore.Filter(ctx, query)
}

func prepareUserQuery() (sq.SelectBuilder, func(*sql.Row) (*User, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserCreationDateCol.identifier(),
			UserChangeDateCol.identifier(),
			UserResourceOwnerCol.identifier(),
			UserSequenceCol.identifier(),
			UserStateCol.identifier(),
			UserUsernameCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanNickNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			HumanPreferredLanguageCol.identifier(),
			HumanGenderCol.identifier(),
			HumanAvaterURLCol.identifier(),
			HumanEmailCol.identifier(),
			HumanIsEmailVerifiedCol.identifier(),
			HumanPhoneCol.identifier(),
			HumanIsPhoneVerifiedCol.identifier(),
			MachineNameCol.identifier(),
			MachineDescriptionCol.identifier()).
			From(projectRolesTable.identifier()).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			LeftJoin(join(MachineUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*User, error) {
			u := new(user)
			firstName := sql.NullString{}
			lastName := sql.NullString{}
			nickName := sql.NullString{}
			displayName := sql.NullString{}
			preferredLanguage := sql.NullString{}
			gender := sql.NullInt16{}
			avatarKey := sql.NullString{}
			email := sql.NullString{}
			isEmailVerified := sql.NullBool{}
			phone := sql.NullString{}
			isPhoneVerified := sql.NullBool{}
			name := sql.NullString{}
			description := sql.NullString{}
			err := row.Scan(
				&u.ID,
				&u.CreationDate,
				&u.ChangeDate,
				&u.ResourceOwner,
				&u.Sequence,
				&u.State,
				&u.Username,
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
				&name,
				&description,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Mf0wf", "Errors.User.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-M00sf", "Errors.Internal")
			}
			if firstName != "" {

			}
			return u, nil
		}
}

func prepareUsersQuery() (sq.SelectBuilder, func(*sql.Rows) (*Users, error)) {
	return sq.Select(
			UserIDCol.identifier(),
			UserColumnCreationDate.identifier(),
			UserColumnChangeDate.identifier(),
			UserColumnResourceOwner.identifier(),
			UserColumnSequence.identifier(),
			UserColumnKey.identifier(),
			UserColumnDisplayName.identifier(),
			UserColumnGroupName.identifier()).
			From(projectRolesTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Rows) (*User, error) {
			p := new(Users)
			err := row.Scan(
				&p.ProjectID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
				&p.Key,
				&p.DisplayName,
				&p.Group,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Mf0wf", "Errors.User.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-M00sf", "Errors.Internal")
			}
			return p, nil
		}
}
