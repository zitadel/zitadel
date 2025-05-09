package stmt

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userStatement struct {
	statement[domain.User]
}

func User(client database.QueryExecutor) *userStatement {
	return &userStatement{
		statement: statement[domain.User]{
			schema: "zitadel",
			table:  "users",
			alias:  "u",
			client: client,
			columns: []Column[domain.User]{
				userColumns[UserInstanceID],
				userColumns[UserOrgID],
				userColumns[UserColumnID],
				userColumns[UserColumnUsername],
				userColumns[UserCreatedAt],
				userColumns[UserUpdatedAt],
				userColumns[UserDeletedAt],
			},
		},
	}
}

func (s *userStatement) Where(condition Condition[domain.User]) *userStatement {
	s.condition = condition
	return s
}

func (s *userStatement) Limit(limit uint32) *userStatement {
	s.limit = limit
	return s
}

func (s *userStatement) Offset(offset uint32) *userStatement {
	s.offset = offset
	return s
}

func (s *userStatement) Get(ctx context.Context) (*domain.User, error) {
	var user domain.User
	err := s.client.QueryRow(ctx, s.query(), s.statement.args...).Scan(s.scanners(&user)...)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userStatement) List(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	rows, err := s.client.Query(ctx, s.query(), s.statement.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err = rows.Scan(s.scanners(&user)...)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (s *userStatement) SetUsername(ctx context.Context, username string) error {
	return nil
}

type UserColumn uint8

var (
	userColumns map[UserColumn]Column[domain.User] = map[UserColumn]Column[domain.User]{
		UserInstanceID: columnDescriptor[domain.User]{
			name: "instance_id",
			scan: func(u *domain.User) any {
				return &u.InstanceID
			},
		},
		UserOrgID: columnDescriptor[domain.User]{
			name: "org_id",
			scan: func(u *domain.User) any {
				return &u.OrgID
			},
		},
		UserColumnID: columnDescriptor[domain.User]{
			name: "id",
			scan: func(u *domain.User) any {
				return &u.ID
			},
		},
		UserColumnUsername: ignoreCaseColumnDescriptor[domain.User]{
			columnDescriptor: columnDescriptor[domain.User]{
				name: "username",
				scan: func(u *domain.User) any {
					return &u.Username
				},
			},
			fieldNameSuffix: "_lower",
		},
		UserCreatedAt: columnDescriptor[domain.User]{
			name: "created_at",
			scan: func(u *domain.User) any {
				return &u.CreatedAt
			},
		},
		UserUpdatedAt: columnDescriptor[domain.User]{
			name: "updated_at",
			scan: func(u *domain.User) any {
				return &u.UpdatedAt
			},
		},
		UserDeletedAt: columnDescriptor[domain.User]{
			name: "deleted_at",
			scan: func(u *domain.User) any {
				return &u.DeletedAt
			},
		},
	}
	humanColumns = map[UserColumn]Column[domain.User]{
		UserHumanColumnEmail: ignoreCaseColumnDescriptor[domain.User]{
			columnDescriptor: columnDescriptor[domain.User]{
				name: "email",
				scan: func(u *domain.User) any {
					human, ok := u.Traits.(*domain.Human)
					if !ok {
						return nil
					}
					if human.Email == nil {
						human.Email = new(domain.Email)
					}
					return &human.Email.Address
				},
			},
			fieldNameSuffix: "_lower",
		},
		UserHumanColumnEmailVerified: columnDescriptor[domain.User]{
			name: "email_is_verified",
			scan: func(u *domain.User) any {
				human, ok := u.Traits.(*domain.Human)
				if !ok {
					return nil
				}
				if human.Email == nil {
					human.Email = new(domain.Email)
				}
				return &human.Email.IsVerified
			},
		},
	}
	machineColumns = map[UserColumn]Column[domain.User]{
		UserMachineDescription: columnDescriptor[domain.User]{
			name: "description",
			scan: func(u *domain.User) any {
				machine, ok := u.Traits.(*domain.Machine)
				if !ok {
					return nil
				}
				if machine == nil {
					machine = new(domain.Machine)
				}
				return &machine.Description
			},
		},
	}
)

const (
	UserInstanceID UserColumn = iota + 1
	UserOrgID
	UserColumnID
	UserColumnUsername
	UserHumanColumnEmail
	UserHumanColumnEmailVerified
	UserMachineDescription
	UserCreatedAt
	UserUpdatedAt
	UserDeletedAt
)
