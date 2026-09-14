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
	// Owners is a bag of kind:id tags used for bulk lifecycle deletes.
	Owners []string
	// OwnerKind is set for UniqueConstraintRemoveByOwner.
	OwnerKind string
	// OwnerID is set for UniqueConstraintRemoveByOwner.
	OwnerID string
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

// OwnerTag returns kind:id. Empty kind or id yields an empty string so it is skipped by WithOwners.
func OwnerTag(kind, id string) string {
	if kind == "" || id == "" {
		return ""
	}
	return kind + ":" + id
}

// WithOwners stores non-empty owner tags on the constraint and returns the receiver for chaining.
func (u *UniqueConstraint) WithOwners(tags ...string) *UniqueConstraint {
	if u == nil {
		return nil
	}
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
	return &UniqueConstraint{
		Action:    UniqueConstraintRemoveByOwner,
		OwnerKind: kind,
		OwnerID:   id,
	}
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
