package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userMachine struct {
	user
}

// AddKey implements domain.MachineUserRepository.
func (u *userMachine) AddKey(key *domain.MachineKey) database.Change {
	panic("unimplemented")
}

// AddMetadata implements domain.MachineUserRepository.
func (u *userMachine) AddMetadata(metadata ...*domain.Metadata) database.Change {
	panic("unimplemented")
}

// AddPersonalAccessToken implements domain.MachineUserRepository.
func (u *userMachine) AddPersonalAccessToken(pat *domain.PersonalAccessToken) database.Change {
	panic("unimplemented")
}

// CreatedAtColumn implements domain.MachineUserRepository.
func (u *userMachine) CreatedAtColumn() database.Column {
	panic("unimplemented")
}

// ExistsMetadata implements domain.MachineUserRepository.
func (u *userMachine) ExistsMetadata(condition database.Condition) database.Condition {
	panic("unimplemented")
}

// IDCondition implements domain.MachineUserRepository.
func (u *userMachine) IDCondition(userID string) database.Condition {
	panic("unimplemented")
}

// InstanceIDCondition implements domain.MachineUserRepository.
func (u *userMachine) InstanceIDCondition(instanceID string) database.Condition {
	panic("unimplemented")
}

// LoginNameCondition implements domain.MachineUserRepository.
func (u *userMachine) LoginNameCondition(op database.TextOperation, loginName string) database.Condition {
	panic("unimplemented")
}

// MetadataKeyCondition implements domain.MachineUserRepository.
func (u *userMachine) MetadataKeyCondition(op database.TextOperation, key string) database.Condition {
	panic("unimplemented")
}

// MetadataValueCondition implements domain.MachineUserRepository.
func (u *userMachine) MetadataValueCondition(op database.BytesOperation, value []byte) database.Condition {
	panic("unimplemented")
}

// OrganizationIDCondition implements domain.MachineUserRepository.
func (u *userMachine) OrganizationIDCondition(orgID string) database.Condition {
	panic("unimplemented")
}

// PrimaryKeyColumns implements domain.MachineUserRepository.
func (u *userMachine) PrimaryKeyColumns() []database.Column {
	panic("unimplemented")
}

// PrimaryKeyCondition implements domain.MachineUserRepository.
func (u *userMachine) PrimaryKeyCondition(instanceID string, userID string) database.Condition {
	panic("unimplemented")
}

// RemoveKey implements domain.MachineUserRepository.
func (u *userMachine) RemoveKey(id string) database.Change {
	panic("unimplemented")
}

// RemoveMetadata implements domain.MachineUserRepository.
func (u *userMachine) RemoveMetadata(condition database.Condition) database.Change {
	panic("unimplemented")
}

// RemovePersonalAccessToken implements domain.MachineUserRepository.
func (u *userMachine) RemovePersonalAccessToken(id string) database.Change {
	panic("unimplemented")
}

// SetAccessTokenType implements domain.MachineUserRepository.
func (u *userMachine) SetAccessTokenType(tokenType domain.PersonalAccessTokenType) database.Change {
	panic("unimplemented")
}

// SetDescription implements domain.MachineUserRepository.
func (u *userMachine) SetDescription(description string) database.Change {
	panic("unimplemented")
}

// SetName implements domain.MachineUserRepository.
func (u *userMachine) SetName(name string) database.Change {
	panic("unimplemented")
}

// SetSecret implements domain.MachineUserRepository.
func (u *userMachine) SetSecret(secret *string) database.Change {
	panic("unimplemented")
}

// SetState implements domain.MachineUserRepository.
func (u *userMachine) SetState(state domain.UserState) database.Change {
	panic("unimplemented")
}

// SetUpdatedAt implements domain.MachineUserRepository.
func (u *userMachine) SetUpdatedAt(updatedAt time.Time) database.Change {
	panic("unimplemented")
}

// SetUsername implements domain.MachineUserRepository.
func (u *userMachine) SetUsername(username string) database.Change {
	panic("unimplemented")
}

// StateColumn implements domain.MachineUserRepository.
func (u *userMachine) StateColumn() database.Column {
	panic("unimplemented")
}

// StateCondition implements domain.MachineUserRepository.
func (u *userMachine) StateCondition(state domain.UserState) database.Condition {
	panic("unimplemented")
}

// TypeColumn implements domain.MachineUserRepository.
func (u *userMachine) TypeColumn() database.Column {
	panic("unimplemented")
}

// TypeCondition implements domain.MachineUserRepository.
func (u *userMachine) TypeCondition(userType domain.UserType) database.Condition {
	panic("unimplemented")
}

// Update implements domain.MachineUserRepository.
func (u *userMachine) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	panic("unimplemented")
}

// UsernameColumn implements domain.MachineUserRepository.
func (u *userMachine) UsernameColumn() database.Column {
	panic("unimplemented")
}

// UsernameCondition implements domain.MachineUserRepository.
func (u *userMachine) UsernameCondition(op database.TextOperation, username string) database.Condition {
	panic("unimplemented")
}

var _ domain.MachineUserRepository = (*userMachine)(nil)
