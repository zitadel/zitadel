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
)

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
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		owners = append(owners, tag)
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
