package projection

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestDomainsProjection_reduceOrgDomainAdded(t *testing.T) {
	projection := &domainsProjection{}

	event := &org.DomainAddedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "org-id",
				InstanceID:    "instance-id",
				Type:          org.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
		Domain: "test.example.com",
	}

	stmt, err := projection.reduceOrgDomainAdded(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	createStmt, ok := stmt.(*handler.CreateStatement)
	require.True(t, ok)

	// Verify the columns being set
	expectedColumns := map[string]interface{}{
		DomainsInstanceIDCol:     "instance-id",
		DomainsOrgIDCol:          "org-id",
		DomainsDomainCol:         "test.example.com",
		DomainsIsVerifiedCol:     false,
		DomainsIsPrimaryCol:      false,
		DomainsValidationTypeCol: domain.OrgDomainValidationTypeUnspecified,
		DomainsCreatedAtCol:      event.CreationDate(),
		DomainsUpdatedAtCol:      event.CreationDate(),
	}

	assert.Len(t, createStmt.Cols, len(expectedColumns))

	for i, col := range createStmt.Cols {
		switch col.Name {
		case DomainsInstanceIDCol:
			assert.Equal(t, "instance-id", col.Value)
		case DomainsOrgIDCol:
			assert.Equal(t, "org-id", col.Value)
		case DomainsDomainCol:
			assert.Equal(t, "test.example.com", col.Value)
		case DomainsIsVerifiedCol:
			assert.Equal(t, false, col.Value)
		case DomainsIsPrimaryCol:
			assert.Equal(t, false, col.Value)
		case DomainsValidationTypeCol:
			assert.Equal(t, domain.OrgDomainValidationTypeUnspecified, col.Value)
		case DomainsCreatedAtCol, DomainsUpdatedAtCol:
			assert.Equal(t, event.CreationDate(), col.Value)
		default:
			t.Errorf("Unexpected column: %s at index %d", col.Name, i)
		}
	}
}

func TestDomainsProjection_reduceOrgDomainVerified(t *testing.T) {
	projection := &domainsProjection{}

	event := &org.DomainVerifiedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "org-id",
				InstanceID:    "instance-id",
				Type:          org.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
		Domain: "test.example.com",
	}

	stmt, err := projection.reduceOrgDomainVerified(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	updateStmt, ok := stmt.(*handler.UpdateStatement)
	require.True(t, ok)

	// Verify update columns
	assert.Len(t, updateStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, updateStmt.Cols[0].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[0].Value)
	assert.Equal(t, DomainsIsVerifiedCol, updateStmt.Cols[1].Name)
	assert.Equal(t, true, updateStmt.Cols[1].Value)

	// Verify conditions
	assert.Len(t, updateStmt.Conditions, 4)
	
	conditionMap := make(map[string]interface{})
	for _, cond := range updateStmt.Conditions {
		conditionMap[cond.Name] = cond.Value
	}

	assert.Equal(t, "instance-id", conditionMap[DomainsInstanceIDCol])
	assert.Equal(t, "org-id", conditionMap[DomainsOrgIDCol])
	assert.Equal(t, "test.example.com", conditionMap[DomainsDomainCol])
	assert.Nil(t, conditionMap[DomainsDeletedAtCol])
}

func TestDomainsProjection_reduceOrgPrimaryDomainSet(t *testing.T) {
	projection := &domainsProjection{}

	event := &org.DomainPrimarySetEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "org-id",
				InstanceID:    "instance-id",
				Type:          org.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
		Domain: "test.example.com",
	}

	stmt, err := projection.reduceOrgPrimaryDomainSet(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	multiStmt, ok := stmt.(*handler.MultiStatement)
	require.True(t, ok)

	// Should have 2 update statements: unset old primary, set new primary
	assert.Len(t, multiStmt.Statements, 2)

	// First statement: unset existing primary
	unsetStmt, ok := multiStmt.Statements[0].(*handler.UpdateStatement)
	require.True(t, ok)
	
	assert.Len(t, unsetStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, unsetStmt.Cols[0].Name)
	assert.Equal(t, DomainsIsPrimaryCol, unsetStmt.Cols[1].Name)
	assert.Equal(t, false, unsetStmt.Cols[1].Value)

	// Second statement: set new primary
	setStmt, ok := multiStmt.Statements[1].(*handler.UpdateStatement)
	require.True(t, ok)
	
	assert.Len(t, setStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, setStmt.Cols[0].Name)
	assert.Equal(t, DomainsIsPrimaryCol, setStmt.Cols[1].Name)
	assert.Equal(t, true, setStmt.Cols[1].Value)
}

func TestDomainsProjection_reduceInstanceDomainAdded(t *testing.T) {
	projection := &domainsProjection{}

	event := &instance.DomainAddedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "instance-id",
				InstanceID:    "instance-id",
				Type:          instance.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
		Domain:    "instance.example.com",
		Generated: false,
	}

	stmt, err := projection.reduceInstanceDomainAdded(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	createStmt, ok := stmt.(*handler.CreateStatement)
	require.True(t, ok)

	// Verify the columns being set for instance domain
	expectedColumns := map[string]interface{}{
		DomainsInstanceIDCol:     "instance-id",
		DomainsOrgIDCol:          nil, // Instance domains have no org_id
		DomainsDomainCol:         "instance.example.com",
		DomainsIsVerifiedCol:     true, // Instance domains are always verified
		DomainsIsPrimaryCol:      false,
		DomainsValidationTypeCol: nil, // Instance domains have no validation type
		DomainsCreatedAtCol:      event.CreationDate(),
		DomainsUpdatedAtCol:      event.CreationDate(),
	}

	assert.Len(t, createStmt.Cols, len(expectedColumns))

	for _, col := range createStmt.Cols {
		switch col.Name {
		case DomainsInstanceIDCol:
			assert.Equal(t, "instance-id", col.Value)
		case DomainsOrgIDCol:
			assert.Nil(t, col.Value)
		case DomainsDomainCol:
			assert.Equal(t, "instance.example.com", col.Value)
		case DomainsIsVerifiedCol:
			assert.Equal(t, true, col.Value)
		case DomainsIsPrimaryCol:
			assert.Equal(t, false, col.Value)
		case DomainsValidationTypeCol:
			assert.Nil(t, col.Value)
		case DomainsCreatedAtCol, DomainsUpdatedAtCol:
			assert.Equal(t, event.CreationDate(), col.Value)
		default:
			t.Errorf("Unexpected column: %s", col.Name)
		}
	}
}

func TestDomainsProjection_reduceOrgDomainRemoved(t *testing.T) {
	projection := &domainsProjection{}

	event := &org.DomainRemovedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "org-id",
				InstanceID:    "instance-id",
				Type:          org.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
		Domain: "test.example.com",
	}

	stmt, err := projection.reduceOrgDomainRemoved(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	updateStmt, ok := stmt.(*handler.UpdateStatement)
	require.True(t, ok)

	// Should update updated_at and deleted_at
	assert.Len(t, updateStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, updateStmt.Cols[0].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[0].Value)
	assert.Equal(t, DomainsDeletedAtCol, updateStmt.Cols[1].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[1].Value)

	// Verify conditions include instance, org, and domain
	assert.Len(t, updateStmt.Conditions, 3)
	
	conditionMap := make(map[string]interface{})
	for _, cond := range updateStmt.Conditions {
		conditionMap[cond.Name] = cond.Value
	}

	assert.Equal(t, "instance-id", conditionMap[DomainsInstanceIDCol])
	assert.Equal(t, "org-id", conditionMap[DomainsOrgIDCol])
	assert.Equal(t, "test.example.com", conditionMap[DomainsDomainCol])
}

func TestDomainsProjection_reduceOrgRemoved(t *testing.T) {
	projection := &domainsProjection{}

	event := &org.OrgRemovedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "org-id",
				InstanceID:    "instance-id",
				Type:          org.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
	}

	stmt, err := projection.reduceOrgRemoved(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	updateStmt, ok := stmt.(*handler.UpdateStatement)
	require.True(t, ok)

	// Should update updated_at and deleted_at
	assert.Len(t, updateStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, updateStmt.Cols[0].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[0].Value)
	assert.Equal(t, DomainsDeletedAtCol, updateStmt.Cols[1].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[1].Value)

	// Should soft delete all domains for the org
	assert.Len(t, updateStmt.Conditions, 3)
	
	conditionMap := make(map[string]interface{})
	for _, cond := range updateStmt.Conditions {
		conditionMap[cond.Name] = cond.Value
	}

	assert.Equal(t, "instance-id", conditionMap[DomainsInstanceIDCol])
	assert.Equal(t, "org-id", conditionMap[DomainsOrgIDCol])
	assert.Nil(t, conditionMap[DomainsDeletedAtCol]) // Only affect non-deleted domains
}

func TestDomainsProjection_reduceInstanceRemoved(t *testing.T) {
	projection := &domainsProjection{}

	event := &instance.InstanceRemovedEvent{
		BaseEvent: &eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{
				ID:            "instance-id",
				InstanceID:    "instance-id",
				Type:          instance.AggregateType,
				ResourceOwner: "instance-id",
			},
			CreationDate: time.Now(),
		},
	}

	stmt, err := projection.reduceInstanceRemoved(event)

	require.NoError(t, err)
	assert.NotNil(t, stmt)

	updateStmt, ok := stmt.(*handler.UpdateStatement)
	require.True(t, ok)

	// Should update updated_at and deleted_at
	assert.Len(t, updateStmt.Cols, 2)
	assert.Equal(t, DomainsUpdatedAtCol, updateStmt.Cols[0].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[0].Value)
	assert.Equal(t, DomainsDeletedAtCol, updateStmt.Cols[1].Name)
	assert.Equal(t, event.CreationDate(), updateStmt.Cols[1].Value)

	// Should soft delete all domains for the instance
	assert.Len(t, updateStmt.Conditions, 2)
	
	conditionMap := make(map[string]interface{})
	for _, cond := range updateStmt.Conditions {
		conditionMap[cond.Name] = cond.Value
	}

	assert.Equal(t, "instance-id", conditionMap[DomainsInstanceIDCol])
	assert.Nil(t, conditionMap[DomainsDeletedAtCol]) // Only affect non-deleted domains
}