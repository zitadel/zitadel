package service

import (
	"context"

	"github.com/zitadel/zitadel/backend/internal/port"
)

type InstanceService struct {
	repo       InstanceRepository
	domainRepo InstanceDomainRepository
	userRepo   UserRepository
	memberRepo InstanceMemberRepository
	inviteRepo InviteInstanceMemberRepository

	domainGenerator DomainGenerator
	idGenerator     IDGenerator

	pool port.Pool[Instance]
}

var _ port.Object = Instance{}

type Instance struct {
	ID      string        `consistent:"id,pk"`
	Name    string        `consistent:"name"`
	State   InstanceState `consistent:"state"`
	Domains []*Domain     `consistent:"domains"`
}

func (i Instance) Columns() []*port.Column {
	return []*port.Column{
		{Name: "id", Value: i.ID},
		{Name: "name", Value: i.Name},
		{Name: "state", Value: i.State},
		{Name: "domains", Value: i.Domains},
	}
}

type InstanceState uint8

const (
	InstanceStateUnspecified InstanceState = iota
	InstanceStateActive
	InstanceStateRemoved
)

// func NewInstance(
// 	repo port.InstanceRepository,
// 	domainGenerator port.DomainGenerator,
// ) *Instance {
// 	return &Instance{repo: repo}
// }

type SetUpInstanceRequest struct {
	Name         string
	CustomDomain *string

	// Admin is the user to be created as the first user of the instance
	// If left empty an invite code will be returned to create the admin user
	Admin *CreateUserRequest
}

type SetUpInstanceResponse struct {
	Instance *Instance
	// Admin is the user that was created as the first user of the instance
	// If the Admin field in the request was empty this field will be nil
	Admin *User
	// InviteCode is the invite code that can be used to create the admin user
	// If the Admin field in the request was not empty this field will be empty
	InviteCode string
}

func (s *InstanceService) SetUpInstance(ctx context.Context, request *SetUpInstanceRequest) (response *SetUpInstanceResponse, err error) {
	instance := &Instance{
		Name:    request.Name,
		State:   InstanceStateActive,
		Domains: make([]*Domain, 0, 2),
	}
	if request.CustomDomain != nil {
		instance.Domains = append(instance.Domains, &Domain{
			Domain:      *request.CustomDomain,
			IsPrimary:   false,
			IsGenerated: false,
		},
		)
	}
	instance.ID, err = s.idGenerator.Generate()
	if err != nil {
		return nil, err
	}
	generatedDomain, err := s.domainGenerator.GenerateDomain()
	if err != nil {
		return nil, err
	}
	instance.Domains = append(instance.Domains, &Domain{
		Domain:      generatedDomain,
		IsPrimary:   true,
		IsGenerated: true,
	})

	var (
		user       *User
		member     *InstanceMember
		inviteCode string
	)

	if request.Admin != nil {
		user = &User{
			Username: request.Admin.Username,
		}
		user.ID, err = s.idGenerator.Generate()
		if err != nil {
			return nil, err
		}
		member = &InstanceMember{
			Member: Member{
				UserID: user.ID,
			},
			Roles: []InstanceMemberRole{InstanceMemberRoleOwner, InstanceMemberRoleAdmin},
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.End(ctx, err)

	if err = s.repo.CreateInstance(ctx, tx, instance); err != nil {
		return nil, err
	}
	if err = s.domainRepo.CreateInstanceDomains(ctx, tx, instance.ID, instance.Domains); err != nil {
		return nil, err
	}

	if user == nil {
		inviteCode, err = s.inviteRepo.InviteInstanceMember(ctx, tx, InstanceMemberRoleOwner, InstanceMemberRoleAdmin)
		if err != nil {
			return nil, err
		}
	} else {
		if err = s.userRepo.CreateUser(ctx, tx, user); err != nil {
			return nil, err
		}
		if err = s.memberRepo.CreateInstanceMember(ctx, tx, member); err != nil {
			return nil, err
		}
	}

	return &SetUpInstanceResponse{
		Instance:   instance,
		Admin:      user,
		InviteCode: inviteCode,
	}, nil
}

type InstanceRepository interface {
	// CreateInstance creates a new instance
	CreateInstance(ctx context.Context, executor port.Executor[Instance], instance *Instance) error
}
