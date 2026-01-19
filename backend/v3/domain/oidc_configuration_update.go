package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// TODO(IAM-Marco): Finish implementation
type OIDCConfigurationUpdate struct {
	instanceDomainName     string
	projectID              string
	managementConsoleAppID string

	// oidcConfigChanges []OIDConfigChange
	// appID string
}

// RequiresTransaction implements [Transactional].
func (o *OIDCConfigurationUpdate) RequiresTransaction() {}

func NewOIDCConfigurationUpdate(instanceDomainName, projectID, consoleAppID string) *OIDCConfigurationUpdate {
	return &OIDCConfigurationUpdate{
		instanceDomainName:     instanceDomainName,
		projectID:              projectID,
		managementConsoleAppID: consoleAppID,
	}
}

// Events implements [Commander].
func (o *OIDCConfigurationUpdate) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	// oidcConfigChangedEvent, err := project.NewOIDCConfigChangedEvent(
	// 	ctx,
	// 	&project.NewAggregate(o.appID, o.projectID).Aggregate,
	// 	o.appID,
	// 	changes,
	// 	o.oidcConfigChanges,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// return []eventstore.Command{
	// 	oidcConfigChangedEvent,
	// }, nil
	return nil, nil
}

// Execute implements [Commander].
func (o *OIDCConfigurationUpdate) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

// String implements [Commander].
func (o *OIDCConfigurationUpdate) String() string {
	return "OIDCConfigurationUpdate"
}

// Validate implements [Commander].
func (o *OIDCConfigurationUpdate) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

var (
	_ (Commander)     = (*OIDCConfigurationUpdate)(nil)
	_ (Transactional) = (*OIDCConfigurationUpdate)(nil)
)
