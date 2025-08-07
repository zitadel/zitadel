package resources

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

type userPatcher struct {
	ctx                  context.Context
	user                 *ScimUser
	metadataChanges      map[metadata.Key]*domain.Metadata
	metadataKeysToRemove map[metadata.Key]bool
	handler              *UsersHandler
}

func (h *UsersHandler) applyPatchesToChangeHuman(ctx context.Context, user *ScimUser, operations patch.OperationCollection) (*command.ChangeHuman, error) {
	patcher := &userPatcher{
		ctx:                  ctx,
		user:                 user,
		metadataChanges:      make(map[metadata.Key]*domain.Metadata),
		metadataKeysToRemove: make(map[metadata.Key]bool),
		handler:              h,
	}

	if err := operations.Apply(patcher, user); err != nil {
		return nil, err
	}

	// we rely on the change detection of the write model to only execute commands that really change data
	changeCommand, err := h.mapToChangeHuman(ctx, user)
	if err != nil {
		return nil, err
	}

	patcher.applyMetadataChangesToCommand(changeCommand)
	return changeCommand, nil
}

func (p *userPatcher) FilterEvaluator() *filter.Evaluator {
	return p.handler.filterEvaluator
}

func (p *userPatcher) Added(attributePath []string) error {
	return p.updateMetadata(attributePath)
}

func (p *userPatcher) Replaced(attributePath []string) error {
	return p.updateMetadata(attributePath)
}

func (p *userPatcher) Removed(attributePath []string) error {
	return p.updateMetadata(attributePath)
}

func (p *userPatcher) applyMetadataChangesToCommand(command *command.ChangeHuman) {
	command.MetadataKeysToRemove = make([]string, 0, len(p.metadataKeysToRemove))
	for key := range p.metadataKeysToRemove {
		command.MetadataKeysToRemove = append(command.MetadataKeysToRemove, string(key))
	}

	command.Metadata = make([]*domain.Metadata, 0, len(p.metadataChanges))
	for _, update := range p.metadataChanges {
		command.Metadata = append(command.Metadata, update)
	}
}

func (p *userPatcher) updateMetadata(attributePath []string) error {
	if len(attributePath) == 0 {
		return nil
	}

	// try full path first (e.g. name.middleName)
	// try root only if full path did not match (e.g. for entitlements.value only entitlements is mapped)
	var ok bool
	var keys []metadata.Key
	if len(attributePath) > 1 {
		keys, ok = metadata.AttributePathToMetadataKeys[strings.Join(attributePath, ".")]
	}

	if !ok {
		keys, ok = metadata.AttributePathToMetadataKeys[attributePath[0]]
		if !ok {
			return nil
		}
	}

	for _, key := range keys {
		value, err := getValueForMetadataKey(p.user, key)
		if err != nil {
			return err
		}

		if len(value) > 0 {
			delete(p.metadataKeysToRemove, key)
			p.metadataChanges[key] = &domain.Metadata{
				Key:   string(metadata.ScopeKey(p.ctx, key)),
				Value: value,
			}
		} else {
			p.metadataKeysToRemove[key] = true
			delete(p.metadataChanges, key)
		}
	}
	return nil
}
