package eventstore

import (
	"context"
	"regexp"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
)

// Aggregate is the basic implementation of Aggregater
type Aggregate struct {
	// ID is the unique identitfier of this aggregate
	ID string `json:"-"`
	// Type is the name of the aggregate.
	Type AggregateType `json:"-"`
	// ResourceOwner is the org this aggregates belongs to
	ResourceOwner string `json:"-"`
	// InstanceID is the instance this aggregate belongs to
	InstanceID string `json:"-"`
	// Version is the semver this aggregate represents
	Version Version `json:"-"`
}

// Version is the old revision of an aggregate
// TODO(adlerhurst): replace version with event.Revision
type Version string

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}

// AggregateOpt is currently public because of the migration of evenstore v2 to v3
type AggregateOpt func(*Aggregate)

// NewAggregate is the default constructor of an aggregate
// opts overwrite values calculated by given parameters
func NewAggregate(
	ctx context.Context,
	id string,
	typ AggregateType,
	version Version,
	opts ...AggregateOpt,
) *Aggregate {
	a := &Aggregate{
		ID:            id,
		Type:          typ,
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		InstanceID:    authz.GetInstance(ctx).InstanceID(),
		Version:       version,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// WithResourceOwner overwrites the resource owner of the aggregate
// by default the resource owner is set by the context
func WithResourceOwner(resourceOwner string) AggregateOpt {
	return func(aggregate *Aggregate) {
		aggregate.ResourceOwner = resourceOwner
	}
}

// AggregateFromWriteModel maps the given WriteModel to an Aggregate
// func AggregateFromWriteModel(
// 	wm *WriteModel,
// 	typ AggregateType,
// 	version Version,
// ) *Aggregate {
// 	return &Aggregate{
// 		ID:            wm.AggregateID,
// 		Type:          typ,
// 		ResourceOwner: wm.ResourceOwner,
// 		InstanceID:    wm.InstanceID,
// 		Version:       version,
// 	}
// }
