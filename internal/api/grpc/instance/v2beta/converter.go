package instance

import (
	"strings"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

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

	queries, err := instanceQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceColumn(req.GetSortingColumn()),
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
	case instance.FieldName_FIELD_NAME_UNSPECIFIED:
		fallthrough
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
	switch q := searchQuery.GetQuery().(type) {
	case *instance.Query_IdQuery:
		return query.NewInstanceIDsListSearchQuery(q.IdQuery.GetIds()...)
	case *instance.Query_DomainQuery:
		return query.NewInstanceDomainsListSearchQuery(q.DomainQuery.GetDomains()...)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-3m0se", "List.Query.Invalid")
	}
}

func CreateInstancePbToSetupInstance(req *instance.CreateInstanceRequest, defaultInstance command.InstanceSetup, externalDomain string) *command.InstanceSetup {
	instance := defaultInstance
	if trimmed := strings.TrimSpace(req.GetInstanceName()); trimmed != "" {
		instance.InstanceName = trimmed
		instance.Org.Name = trimmed
	}
	if trimmed := strings.TrimSpace(req.GetCustomDomain()); trimmed != "" {
		instance.CustomDomain = trimmed
	}
	if trimmed := strings.TrimSpace(req.GetFirstOrgName()); trimmed != "" {
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

	if lang := language.Make(strings.TrimSpace(req.GetDefaultLanguage())); !lang.IsRoot() {
		instance.DefaultLanguage = lang
	}

	return &instance
}

func createInstancePbToAddHuman(req *instance.CreateInstanceRequest_Human, defaultHuman command.AddHuman, userLoginMustBeDomain bool, org, externalDomain string) *command.AddHuman {
	user := defaultHuman
	if req.Email != nil {
		user.Email.Address = domain.EmailAddress(strings.TrimSpace(req.GetEmail().GetEmail()))
		user.Email.Verified = req.GetEmail().GetIsEmailVerified()
	}
	if req.GetProfile() != nil {
		if firstName := strings.TrimSpace(req.GetProfile().GetFirstName()); firstName != "" {
			user.FirstName = firstName
		}
		if lastName := strings.TrimSpace(req.GetProfile().GetLastName()); lastName != "" {
			user.LastName = lastName
		}
		if lang := strings.TrimSpace(req.GetProfile().GetPreferredLanguage()); lang != "" {
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
	if username := strings.TrimSpace(req.GetUserName()); username != "" {
		user.Username = username
	}

	if req.GetPassword() != nil {
		user.Password = req.GetPassword().GetPassword()
		user.PasswordChangeRequired = req.GetPassword().GetPasswordChangeRequired()
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

	if username := strings.TrimSpace(req.GetUserName()); username != "" {
		machine.Machine.Username = username
	}
	if name := strings.TrimSpace(req.GetName()); name != "" {
		machine.Machine.Name = name
	}

	if defaultMachine.Pat != nil || req.GetPersonalAccessToken() != nil {
		pat := command.AddPat{
			// Scopes are currently static and can not be overwritten
			Scopes: []string{oidc.ScopeOpenID, oidc.ScopeProfile, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner},
		}
		if req.GetPersonalAccessToken().GetExpirationDate().IsValid() {
			pat.ExpirationDate = req.GetPersonalAccessToken().GetExpirationDate().AsTime()
		} else if defaultMachine.Pat != nil && !defaultMachine.Pat.ExpirationDate.IsZero() {
			pat.ExpirationDate = defaultMachine.Pat.ExpirationDate
		}
		machine.Pat = &pat
	}

	if defaultMachine.MachineKey != nil || req.GetMachineKey() != nil {
		machineKey := command.AddMachineKey{}
		if defaultMachine.MachineKey != nil {
			machineKey = *defaultMachine.MachineKey
		}
		if req.GetMachineKey() != nil {
			if req.GetMachineKey().GetType() != 0 {
				machineKey.Type = authn.KeyTypeToDomain(req.GetMachineKey().GetType())
			}
			if req.GetMachineKey().ExpirationDate.IsValid() {
				machineKey.ExpirationDate = req.GetMachineKey().GetExpirationDate().AsTime()
			}
		}
		machine.MachineKey = &machineKey
	}

	return machine
}

func ListCustomDomainsRequestToModel(req *instance.ListCustomDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := domainQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceDomainColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func fieldNameToInstanceDomainColumn(fieldName instance.DomainFieldName) query.Column {
	switch fieldName {
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceDomainDomainCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_GENERATED:
		return query.InstanceDomainIsGeneratedCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_PRIMARY:
		return query.InstanceDomainIsPrimaryCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceDomainCreationDateCol
	case instance.DomainFieldName_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return query.Column{}
	}
}

func domainQueriesToModel(queries []*instance.DomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = domainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func domainQueryToModel(searchQuery *instance.DomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.GetQuery().(type) {
	case *instance.DomainSearchQuery_DomainQuery:
		return query.NewInstanceDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.GetMethod()), q.DomainQuery.GetDomain())
	case *instance.DomainSearchQuery_GeneratedQuery:
		return query.NewInstanceDomainGeneratedSearchQuery(q.GeneratedQuery.GetGenerated())
	case *instance.DomainSearchQuery_PrimaryQuery:
		return query.NewInstanceDomainPrimarySearchQuery(q.PrimaryQuery.GetPrimary())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func ListTrustedDomainsRequestToModel(req *instance.ListTrustedDomainsRequest, defaults systemdefaults.SystemDefaults) (*query.InstanceTrustedDomainSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(defaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := trustedDomainQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}

	return &query.InstanceTrustedDomainSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToInstanceTrustedDomainColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func trustedDomainQueriesToModel(queries []*instance.TrustedDomainSearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = trustedDomainQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func trustedDomainQueryToModel(searchQuery *instance.TrustedDomainSearchQuery) (query.SearchQuery, error) {
	switch q := searchQuery.GetQuery().(type) {
	case *instance.TrustedDomainSearchQuery_DomainQuery:
		return query.NewInstanceTrustedDomainDomainSearchQuery(object.TextMethodToQuery(q.DomainQuery.GetMethod()), q.DomainQuery.GetDomain())
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ags42", "List.Query.Invalid")
	}
}

func trustedDomainsToPb(domains []*query.InstanceTrustedDomain) []*instance.TrustedDomain {
	d := make([]*instance.TrustedDomain, len(domains))
	for i, domain := range domains {
		d[i] = trustedDomainToPb(domain)
	}
	return d
}

func trustedDomainToPb(d *query.InstanceTrustedDomain) *instance.TrustedDomain {
	return &instance.TrustedDomain{
		Domain: d.Domain,
		Details: object.ToViewDetailsPb(
			d.Sequence,
			d.CreationDate,
			d.ChangeDate,
			d.InstanceID,
		),
	}
}

func fieldNameToInstanceTrustedDomainColumn(fieldName instance.TrustedDomainFieldName) query.Column {
	switch fieldName {
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_DOMAIN:
		return query.InstanceTrustedDomainDomainCol
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_CREATION_DATE:
		return query.InstanceTrustedDomainCreationDateCol
	case instance.TrustedDomainFieldName_TRUSTED_DOMAIN_FIELD_NAME_UNSPECIFIED:
		fallthrough
	default:
		return query.Column{}
	}
}
