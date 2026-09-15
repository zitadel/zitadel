package eventstore

type UniqueConstraint struct {
	// UniqueType is the table name for the unique constraint
	UniqueType string
	// UniqueField is the unique key
	UniqueField string
	// Action defines if unique constraint should be added or removed
	Action UniqueConstraintAction
	// ErrorMessage defines the translation file key for the error message
	ErrorMessage string
	// IsGlobal defines if the unique constraint is globally unique or just within a single instance
	IsGlobal bool
	// Owners is a bag of kind:id tags used for bulk lifecycle deletes
	Owners []string
}

type UniqueConstraintAction int8

const (
	UniqueConstraintAdd UniqueConstraintAction = iota
	UniqueConstraintRemove
	UniqueConstraintInstanceRemove
	UniqueConstraintRemoveByOwner

	uniqueConstraintActionCount
)

const (
	UniqueConstraintOwnerOrg     = "org"
	UniqueConstraintOwnerUser    = "user"
	UniqueConstraintOwnerIDP     = "idp"
	UniqueConstraintOwnerProject = "project"
	UniqueConstraintOwnerGrant   = "grant"

	// UniqueConstraintOwnersBackfillStep is the setup migration that stamps owners from projections.
	UniqueConstraintOwnersBackfillStep = "79_backfill_unique_constraint_owners"
)

// UniqueTypesWithOwners are unique_types that live writes and setup 79 stamp with owner tags.
// Instance mail_text shares unique_type "mail_text" but stays empty by design, so it is omitted.
var UniqueTypesWithOwners = []string{
	"usernames",
	"external_idps",
	"org_name",
	"org_domain",
	"project_names",
	"appname",
	"project_role",
	"entity_ids",
	"project_grant",
	"project_grant_member",
	"user_grant",
	"member",
	"group_name",
	"action_names",
	"idp_config_names",
}

func (f UniqueConstraintAction) Valid() bool {
	return f >= 0 && f < uniqueConstraintActionCount
}

func OwnerTag(kind, id string) string {
	if kind == "" || id == "" {
		return ""
	}
	return kind + ":" + id
}

func (u *UniqueConstraint) WithOwners(tags ...string) *UniqueConstraint {
	owners := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag != "" {
			owners = append(owners, tag)
		}
	}
	u.Owners = owners
	return u
}

func NewAddEventUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		Action:       UniqueConstraintAdd,
		Owners:       []string{},
	}
}

func NewRemoveUniqueConstraint(
	uniqueType,
	uniqueField string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		Action:      UniqueConstraintRemove,
	}
}

func NewRemoveInstanceUniqueConstraints() *UniqueConstraint {
	return &UniqueConstraint{
		Action: UniqueConstraintInstanceRemove,
	}
}

func NewRemoveUniqueConstraintsByOwner(kind, id string) *UniqueConstraint {
	constraint := &UniqueConstraint{
		Action: UniqueConstraintRemoveByOwner,
	}
	if tag := OwnerTag(kind, id); tag != "" {
		constraint.Owners = []string{tag}
	}
	return constraint
}

func NewAddGlobalUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		IsGlobal:     true,
		Action:       UniqueConstraintAdd,
		Owners:       []string{},
	}
}

func NewRemoveGlobalUniqueConstraint(
	uniqueType,
	uniqueField string) *UniqueConstraint {
	return &UniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		IsGlobal:    true,
		Action:      UniqueConstraintRemove,
	}
}
