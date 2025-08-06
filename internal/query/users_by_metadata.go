package query

import (
	"context"
	"database/sql"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UsersByMetadataSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

type UserByMetadata struct {
	ResourceOwner string `json:"resource_owner,omitempty"`
	Key           string `json:"key,omitempty"`
	Value         []byte `json:"value,omitempty"`
	User          *User  `json:"user,omitempty"`
}

type UsersByMetadata struct {
	SearchResponse
	UsersByMeta []*UserByMetadata
}

func (q *Queries) SearchUsersByMetadata(ctx context.Context, queries *UsersByMetadataSearchQueries, permissionCheck domain.PermissionCheck) (usersByMeta *UsersByMetadata, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	permissionCheckV2 := PermissionV2(ctx, permissionCheck)

	query, scan := prepareUserByMetadataListQuery()
	query = userPermissionCheckV2(ctx, query, permissionCheckV2, queries.Queries)

	for _, q := range queries.Queries {
		query = q.toQuery(query)
	}

	eq := sq.Eq{
		UserMetadataInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}

	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-J1CvYA", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		usersByMeta, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		logging.Errorf("Got error: %v", err)
		return nil, err
	}

	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		usersByMetadataCheckPermission(ctx, usersByMeta, permissionCheck)
	}

	return usersByMeta, nil
}

func prepareUserByMetadataListQuery() (sq.SelectBuilder, func(*sql.Rows) (*UsersByMetadata, error)) {
	return sq.Select(
			UserMetadataKeyCol.identifier(),
			UserMetadataValueCol.identifier(),

			UserIDCol.identifier(),
			UserStateCol.identifier(),
			UserUsernameCol.identifier(),
			UserTypeCol.identifier(),

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
			HumanMFAInitSkippedCol.identifier(),

			MachineUserIDCol.identifier(),
			MachineNameCol.identifier(),
			MachineDescriptionCol.identifier(),
			MachineSecretCol.identifier(),
			MachineAccessTokenTypeCol.identifier(),

			countColumn.identifier()).
			From(userMetadataTable.identifier()).
			Join(join(UserIDCol, UserMetadataUserIDCol)).
			LeftJoin(join(HumanUserIDCol, UserIDCol)).
			LeftJoin(join(MachineUserIDCol, UserIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*UsersByMetadata, error) {
			usersByMeta := make([]*UserByMetadata, 0)
			var count uint64
			for rows.Next() {
				userByMeta := new(UserByMetadata)
				userByMeta.User = new(User)
				userByMeta.User.Human = new(Human)
				userByMeta.User.Machine = new(Machine)

				human, machine := sqlHuman{}, sqlMachine{}

				err := rows.Scan(
					&userByMeta.Key,
					&userByMeta.Value,

					&userByMeta.User.ID,
					&userByMeta.User.State,
					&userByMeta.User.Username,
					&userByMeta.User.Type,

					&human.humanID,
					&human.firstName,
					&human.lastName,
					&human.nickName,
					&human.displayName,
					&human.preferredLanguage,
					&human.gender,
					&human.avatarKey,
					&human.email,
					&human.isEmailVerified,
					&human.phone,
					&human.isPhoneVerified,
					&human.passwordChangeRequired,
					&human.passwordChanged,
					&human.mfaInitSkipped,

					&machine.machineID,
					&machine.name,
					&machine.description,
					&machine.encodedSecret,
					&machine.accessTokenType,

					&count,
				)
				if err != nil {
					return nil, err
				}

				if human.humanID.Valid {
					userByMeta.User.Human.FirstName = human.firstName.String
					userByMeta.User.Human.LastName = human.lastName.String
					userByMeta.User.Human.NickName = human.nickName.String
					userByMeta.User.Human.DisplayName = human.displayName.String

					userByMeta.User.Human.PreferredLanguage, err = language.Parse(human.preferredLanguage.String)
					if err != nil {
						return nil, err
					}

					userByMeta.User.Human.Gender = domain.Gender(human.gender.Int32)
					userByMeta.User.Human.AvatarKey = human.avatarKey.String
					userByMeta.User.Human.Email = domain.EmailAddress(human.email.String)
					userByMeta.User.Human.IsEmailVerified = human.isEmailVerified.Bool
					userByMeta.User.Human.Phone = domain.PhoneNumber(human.phone.String)
					userByMeta.User.Human.IsPhoneVerified = human.isPhoneVerified.Bool
					userByMeta.User.Human.PasswordChangeRequired = human.passwordChangeRequired.Bool
					userByMeta.User.Human.PasswordChanged = human.passwordChanged.Time
					userByMeta.User.Human.MFAInitSkipped = human.mfaInitSkipped.Time
				}

				if machine.machineID.Valid {
					userByMeta.User.Machine.Name = machine.name.String
					userByMeta.User.Machine.Description = machine.description.String
					userByMeta.User.Machine.EncodedSecret = machine.encodedSecret.String
					userByMeta.User.Machine.AccessTokenType = domain.OIDCTokenType(machine.accessTokenType.Int32)
				}

				usersByMeta = append(usersByMeta, userByMeta)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-hwTK0J", "Errors.Query.CloseRows")
			}

			return &UsersByMetadata{
				UsersByMeta: usersByMeta,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func usersByMetadataCheckPermission(ctx context.Context, users *UsersByMetadata, permissionCheck domain.PermissionCheck) {
	users.UsersByMeta = slices.DeleteFunc(users.UsersByMeta,
		func(user *UserByMetadata) bool {
			return userCheckPermission(ctx, user.ResourceOwner, user.User.ID, permissionCheck) != nil
		},
	)
}
