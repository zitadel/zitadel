package instance

import (
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel/cmd/build"
	authn "github.com/zitadel/zitadel/internal/api/grpc/authn/v2beta"
	filter "github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	instance "github.com/zitadel/zitadel/pkg/grpc/instance/v2beta"
	"golang.org/x/text/language"
)

func InstancesToPb(instances []*query.Instance) []*instance.Instance {
	list := []*instance.Instance{}
	for _, instance := range instances {
		list = append(list, ToProtoObject(instance))
	}
	return list
}

func ToProtoObject(inst *query.Instance) *instance.Instance {
	return &instance.Instance{
		Id:      inst.ID,
		Name:    inst.Name,
		Domains: DomainsToPb(inst.Domains),
		Version: build.Version(),
		Details: object.ToViewDetailsPb(inst.Sequence, inst.CreationDate, inst.ChangeDate, inst.ID),
	}
}

func DomainsToPb(domains []*query.InstanceDomain) []*instance.Domain {
	d := []*instance.Domain{}
	for _, dm := range domains {
		pbDomain := DomainToPb(dm)
		d = append(d, pbDomain)
	}
	return d
}

func DomainToPb(d *query.InstanceDomain) *instance.Domain {
	return &instance.Domain{
		Domain:    d.Domain,
		Primary:   d.IsPrimary,
		Generated: d.IsGenerated,
		Details: object.ToViewDetailsPb(
			d.Sequence,
			d.CreationDate,
			d.ChangeDate,
			d.InstanceID,
		),
	}
}

func ListInstancesRequestToModel(req *instance.ListInstancesRequest, sysDefaults systemdefaults.SystemDefaults) (*query.InstanceSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := instanceQueriesToModel(req.Queries)
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

func fieldNameToInstanceColumn(fieldName instance.FieldName) query.Column {
	switch fieldName {
	case instance.FieldName_FIELD_NAME_ID:
		return query.InstanceColumnID
	case instance.FieldName_FIELD_NAME_NAME:
		return query.InstanceColumnName
	case instance.FieldName_FIELD_NAME_CREATION_DATE:
		return query.InstanceColumnCreationDate
	default:
		return query.Column{}
	}
}

func instanceQueriesToModel(queries []*instance.Query) (_ []query.SearchQuery, err error) {
	q := []query.SearchQuery{}
	for _, query := range queries {
		model, err := instanceQueryToModel(query)
		if err != nil {
			return nil, err
		}
		q = append(q, model)
	}
	return q, nil
}

func instanceQueryToModel(searchQuery *instance.Query) (query.SearchQuery, error) {
	switch q := searchQuery.Query.(type) {
	case *instance.Query_IdQuery:
		return query.NewInstanceIDsListSearchQuery(q.IdQuery.Ids...)
	case *instance.Query_DomainQuery:
		return query.NewInstanceDomainsListSearchQuery(q.DomainQuery.Domains...)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}

func CreateInstancePbToSetupInstance(req *instance.CreateInstanceRequest, defaultInstance command.InstanceSetup, externalDomain string) *command.InstanceSetup {
	instance := defaultInstance
	if trimmed := strings.TrimSpace(req.InstanceName); trimmed != "" {
		instance.InstanceName = trimmed
		instance.Org.Name = trimmed
	}
	if trimmed := strings.TrimSpace(req.CustomDomain); trimmed != "" {
		instance.CustomDomain = trimmed
	}
	if trimmed := strings.TrimSpace(req.FirstOrgName); trimmed != "" {
		instance.Org.Name = trimmed
	}

	if user := req.GetMachine(); user != nil {
		defaultMachine := instance.Org.Machine
		if defaultMachine == nil {
			defaultMachine = &command.AddMachine{}
		}

		instance.Org.Machine = createInstancePbToAddMachine(user, *defaultMachine)
		instance.Org.Human = nil
	} else if user := req.GetHuman(); user != nil {
		defaultHuman := instance.Org.Human
		if instance.Org.Human != nil {
			defaultHuman = &command.AddHuman{}
		}

		instance.Org.Human = createInstancePbToAddHuman(user, *defaultHuman, instance.DomainPolicy.UserLoginMustBeDomain, instance.Org.Name, externalDomain)
		instance.Org.Machine = nil
	}

	if lang := language.Make(strings.TrimSpace(req.DefaultLanguage)); !lang.IsRoot() {
		instance.DefaultLanguage = lang
	}

	return &instance
}

func createInstancePbToAddHuman(req *instance.CreateInstanceRequest_Human, defaultHuman command.AddHuman, userLoginMustBeDomain bool, org, externalDomain string) *command.AddHuman {
	user := defaultHuman
	if req.Email != nil {
		user.Email.Address = domain.EmailAddress(strings.TrimSpace(req.Email.Email))
		user.Email.Verified = req.Email.IsEmailVerified
	}
	if req.Profile != nil {
		if firstName := strings.TrimSpace(req.Profile.FirstName); firstName != "" {
			user.FirstName = firstName
		}
		if lastName := strings.TrimSpace(req.Profile.LastName); lastName != "" {
			user.LastName = lastName
		}
		if lang := strings.TrimSpace(req.Profile.PreferredLanguage); lang != "" {
			lang, err := language.Parse(lang)
			if err == nil {
				user.PreferredLanguage = lang
			}
		}
	}
	// check if default username is email style or else append @<orgname>.<custom-domain>
	// this way we have the same value as before changing `UserLoginMustBeDomain` to false
	if !userLoginMustBeDomain && !strings.Contains(user.Username, "@") {
		orgDomain, _ := domain.NewIAMDomainName(org, externalDomain)
		user.Username = user.Username + "@" + orgDomain
	}
	if username := strings.TrimSpace(req.UserName); username != "" {
		user.Username = username
	}

	if req.Password != nil {
		user.Password = req.Password.Password
		user.PasswordChangeRequired = req.Password.PasswordChangeRequired
	}
	return &user
}

func createInstancePbToAddMachine(req *instance.CreateInstanceRequest_Machine, defaultMachine command.AddMachine) (machine *command.AddMachine) {
	machine = &command.AddMachine{}
	if defaultMachine.Machine != nil {
		machineCopy := *defaultMachine.Machine
		machine.Machine = &machineCopy
	} else {
		machine.Machine = &command.Machine{}
	}

	if username := strings.TrimSpace(req.UserName); username != "" {
		machine.Machine.Username = username
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		machine.Machine.Name = name
	}

	if defaultMachine.Pat != nil || req.PersonalAccessToken != nil {
		pat := command.AddPat{
			// Scopes are currently static and can not be overwritten
			Scopes: []string{oidc.ScopeOpenID, oidc.ScopeProfile, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner},
		}
		if req.GetPersonalAccessToken().GetExpirationDate().IsValid() {
			pat.ExpirationDate = req.PersonalAccessToken.ExpirationDate.AsTime()
		} else if defaultMachine.Pat != nil && !defaultMachine.Pat.ExpirationDate.IsZero() {
			pat.ExpirationDate = defaultMachine.Pat.ExpirationDate
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
