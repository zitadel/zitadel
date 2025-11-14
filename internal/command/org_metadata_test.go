package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetOrgMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx      context.Context
			orgID    string
			metadata *domain.Metadata
		}
	)
	type res struct {
		want *domain.Metadata
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadata: &domain.Metadata{
					Key: "key",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectPush(
						org.NewMetadataSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key",
							[]byte("value"),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: []byte("value"),
					State: domain.MetadataStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetOrgMetadata(tt.args.ctx, tt.args.orgID, tt.args.metadata)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_BulkSetOrgMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			metadataList []*domain.Metadata
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty meta data list, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "org not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key"},
					{Key: "key1"},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectPush(
						org.NewMetadataSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key",
							[]byte("value"),
						),
						org.NewMetadataSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key1",
							[]byte("value1"),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkSetOrgMetadata(tt.args.ctx, tt.args.orgID, tt.args.metadataList...)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_OrgRemoveMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx         context.Context
			orgID       string
			metadataKey string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				metadataKey: "key",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				metadataKey: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "meta data not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				metadataKey: "key",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
					expectPush(
						org.NewMetadataRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				metadataKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveOrgMetadata(tt.args.ctx, tt.args.orgID, tt.args.metadataKey)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_BulkRemoveOrgMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			metadataList []string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty meta data list, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "org not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "remove metadata keys not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key1",
								[]byte("value1"),
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				metadataList: []string{""},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"ZITADEL",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							org.NewMetadataSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"key1",
								[]byte("value1"),
							),
						),
					),
					expectPush(
						org.NewMetadataRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key",
						),
						org.NewMetadataRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key1",
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkRemoveOrgMetadata(tt.args.ctx, tt.args.orgID, tt.args.metadataList...)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}
