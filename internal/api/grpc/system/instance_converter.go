package system

import (
	"strings"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/grpc/authn"
	instance_grpc "github.com/zitadel/zitadel/internal/api/grpc/instance"
	member_grpc "github.com/zitadel/zitadel/internal/api/grpc/member"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	instance_pb "github.com/zitadel/zitadel/pkg/grpc/instance"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func CreateInstancePbToSetupInstance(req *system_pb.CreateInstanceRequest, defaultInstance command.InstanceSetup, externalDomain string) *command.InstanceSetup {
	instance := defaultInstance
	if req.InstanceName != "" {
		instance.InstanceName = req.InstanceName
		instance.Org.Name = req.InstanceName
	}
	if req.CustomDomain != "" {
		instance.CustomDomain = req.CustomDomain
	}
	if req.FirstOrgName != "" {
		instance.Org.Name = req.FirstOrgName
	}

	if user := req.GetMachine(); user != nil {
		defaultMachine := instance.Org.Machine
		if defaultMachine == nil {
			defaultMachine = new(command.AddMachine)
		}

		instance.Org.Machine = createInstancePbToAddMachine(user, *defaultMachine)
		instance.Org.Human = nil
	} else if user := req.GetHuman(); user != nil {
		defaultHuman := instance.Org.Human
		if instance.Org.Human != nil {
			defaultHuman = new(command.AddHuman)
		}

		instance.Org.Human = createInstancePbToAddHuman(user, *defaultHuman, instance.DomainPolicy.UserLoginMustBeDomain, instance.Org.Name, externalDomain)
		instance.Org.Machine = nil
	}

	if lang := language.Make(req.DefaultLanguage); !lang.IsRoot() {
		instance.DefaultLanguage = lang
	}

	return &instance
}

func createInstancePbToAddHuman(req *system_pb.CreateInstanceRequest_Human, defaultHuman command.AddHuman, userLoginMustBeDomain bool, org, externalDomain string) *command.AddHuman {
	user := defaultHuman
	if req.Email != nil {
		user.Email.Address = domain.EmailAddress(req.Email.Email)
		user.Email.Verified = req.Email.IsEmailVerified
	}
	if req.Profile != nil {
		if req.Profile.FirstName != "" {
			user.FirstName = req.Profile.FirstName
		}
		if req.Profile.LastName != "" {
			user.LastName = req.Profile.LastName
		}
		if req.Profile.PreferredLanguage != "" {
			lang, err := language.Parse(req.Profile.PreferredLanguage)
			if err == nil {
				user.PreferredLanguage = lang
			}
		}
	}
	// check if default username is email style or else append @<orgname>.<custom-domain>
	// this way we have the same value as before changing `UserLoginMustBeDomain` to false
	if !userLoginMustBeDomain && !strings.Contains(user.Username, "@") {
		user.Username = user.Username + "@" + domain.NewIAMDomainName(org, externalDomain)
	}
	if req.UserName != "" {
		user.Username = req.UserName
	}

	if req.Password != nil {
		user.Password = req.Password.Password
		user.PasswordChangeRequired = req.Password.PasswordChangeRequired
	}
	return &user
}

func createInstancePbToAddMachine(req *system_pb.CreateInstanceRequest_Machine, defaultMachine command.AddMachine) (machine *command.AddMachine) {
	machine = new(command.AddMachine)
	if defaultMachine.Machine != nil {
		machineCopy := *defaultMachine.Machine
		machine.Machine = &machineCopy
	} else {
		machine.Machine = new(command.Machine)
	}

	if req.UserName != "" {
		machine.Machine.Username = req.UserName
	}
	if req.Name != "" {
		machine.Machine.Name = req.Name
	}

	if defaultMachine.Pat != nil || req.PersonalAccessToken != nil {
		pat := command.AddPat{
			// Scopes are currently static and can not be overwritten
			Scopes: []string{oidc.ScopeOpenID, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner},
		}

		if !defaultMachine.Pat.ExpirationDate.IsZero() {
			pat.ExpirationDate = defaultMachine.Pat.ExpirationDate
		} else if req.PersonalAccessToken.ExpirationDate.IsValid() {
			pat.ExpirationDate = req.PersonalAccessToken.ExpirationDate.AsTime()
		}

		machine.Pat = &pat
	}

	if defaultMachine.MachineKey != nil || req.MachineKey != nil {
		machineKey := command.AddMachineKey{}
		if defaultMachine.MachineKey != nil {
			machineKey = *defaultMachine.MachineKey
		}
		if req.MachineKey != nil {
			if req.MachineKey.Type != 0 {
				machineKey.Type = authn.KeyTypeToDomain(req.MachineKey.Type)
			}
			if req.MachineKey.ExpirationDate.IsValid() {
				machineKey.ExpirationDate = req.MachineKey.ExpirationDate.AsTime()
			}
		}
		machine.MachineKey = &machineKey
	}

	return machine
}

func AddInstancePbToSetupInstance(req *system_pb.AddInstanceRequest, defaultInstance command.InstanceSetup, externalDomain string) *command.InstanceSetup {
	instance := defaultInstance

	if req.InstanceName != "" {
		instance.InstanceName = req.InstanceName
		instance.Org.Name = req.InstanceName
	}
	if req.CustomDomain != "" {
		instance.CustomDomain = req.CustomDomain
	}
	if req.FirstOrgName != "" {
		instance.Org.Name = req.FirstOrgName
	}

	if defaultInstance.Org.Human != nil {
		// used to not overwrite the default human later
		humanCopy := *defaultInstance.Org.Human
		instance.Org.Human = &humanCopy
	} else {
		instance.Org.Human = new(command.AddHuman)
	}
	if req.OwnerEmail.Email != "" {
		instance.Org.Human.Email.Address = domain.EmailAddress(req.OwnerEmail.Email)
		instance.Org.Human.Email.Verified = req.OwnerEmail.IsEmailVerified
	}
	if req.OwnerProfile != nil {
		if req.OwnerProfile.FirstName != "" {
			instance.Org.Human.FirstName = req.OwnerProfile.FirstName
		}
		if req.OwnerProfile.LastName != "" {
			instance.Org.Human.LastName = req.OwnerProfile.LastName
		}
		if req.OwnerProfile.PreferredLanguage != "" {
			lang, err := language.Parse(req.OwnerProfile.PreferredLanguage)
			if err == nil {
				instance.Org.Human.PreferredLanguage = lang
			}
		}
	}
	if req.OwnerUserName != "" {
		instance.Org.Human.Username = req.OwnerUserName
	}
	// check if default username is email style or else append @<orgname>.<custom-domain>
	// this way we have the same value as before changing `UserLoginMustBeDomain` to false
	if !instance.DomainPolicy.UserLoginMustBeDomain && !strings.Contains(instance.Org.Human.Username, "@") {
		instance.Org.Human.Username = instance.Org.Human.Username + "@" + domain.NewIAMDomainName(instance.Org.Name, externalDomain)
	}
	if req.OwnerPassword != nil {
		instance.Org.Human.Password = req.OwnerPassword.Password
		instance.Org.Human.PasswordChangeRequired = req.OwnerPassword.PasswordChangeRequired
	}
	if lang := language.Make(req.DefaultLanguage); lang != language.Und {
		instance.DefaultLanguage = lang
	}

	return &instance
}

func ListInstancesRequestToModel(req *system_pb.ListInstancesRequest) (*query.InstanceSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.InstanceQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceColumn(fieldName instance_pb.FieldName) query.Column {
	switch fieldName {
	case instance_pb.FieldName_FIELD_NAME_ID:
		return query.InstanceColumnID
	case instance_pb.FieldName_FIELD_NAME_NAME:
		return query.InstanceColumnName
	case instance_pb.FieldName_FIELD_NAME_CREATION_DATE:
		return query.InstanceColumnCreationDate
	default:
		return query.Column{}
	}
}

func ListInstanceDomainsRequestToModel(req *system_pb.ListDomainsRequest) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := instance_grpc.DomainQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceDomainColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceDomainColumn(fieldName instance_pb.DomainFieldName) query.Column {
	switch fieldName {
	case instance_pb.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceDomainDomainCol
	case instance_pb.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return query.InstanceDomainIsGeneratedCol
	case instance_pb.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		return query.InstanceDomainIsPrimaryCol
	case instance_pb.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceDomainCreationDateCol
	default:
		return query.Column{}
	}
}

func ListIAMMembersRequestToQuery(req *system_pb.ListIAMMembersRequest) (*query.IAMMembersQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := member_grpc.MemberQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.IAMMembersQuery{
		MembersQuery: query.MembersQuery{
			SearchRequest: query.SearchRequest{
				Offset: offset,
				Limit:  limit,
				Asc:    asc,
				// SortingColumn: model.IAMMemberSearchKey, //TOOD: not implemented in proto
			},
			Queries: queries,
		},
	}, nil
}
