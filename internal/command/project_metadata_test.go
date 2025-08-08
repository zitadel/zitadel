package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetProjectMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx       context.Context
			projectID string
			metadata  *domain.Metadata
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
			name: "project not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPush(
						project.NewMetadataSetEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key",
							[]byte("value"),
						),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "ro-1",
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
			got, err := r.SetProjectMetadata(tt.args.ctx, tt.args.projectID, "ro-1", tt.args.metadata)
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

func TestCommandSide_BulkSetProjectMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			projectID    string
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
				ctx:       context.Background(),
				projectID: "project1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectPush(
						project.NewMetadataSetEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key",
							[]byte("value"),
						),
						project.NewMetadataSetEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key1",
							[]byte("value1"),
						),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ro-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkSetProjectMetadata(tt.args.ctx, tt.args.projectID, "ro-1", tt.args.metadataList...)
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

func TestCommandSide_ProjectRemoveMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx         context.Context
			projectID   string
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
			name: "project not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				projectID:   "project1",
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
				projectID:   "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				projectID:   "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
					expectPush(
						project.NewMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				projectID:   "project1",
				metadataKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ro-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveProjectMetadata(tt.args.ctx, tt.args.projectID, "ro-1", tt.args.metadataKey)
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

func TestCommandSide_BulkRemoveProjectMetadata(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			projectID    string
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
				ctx:       context.Background(),
				projectID: "project1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "project not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:          context.Background(),
				projectID:    "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				projectID:    "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key1",
								[]byte("value1"),
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				projectID:    "project1",
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
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"ZITADEL", true, true, true,
								domain.PrivateLabelingSettingUnspecified,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							project.NewMetadataSetEvent(context.Background(),
								&project.NewAggregate("project1", "ro-1").Aggregate,
								"key1",
								[]byte("value1"),
							),
						),
					),
					expectPush(
						project.NewMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key",
						),
						project.NewMetadataRemovedEvent(context.Background(),
							&project.NewAggregate("project1", "ro-1").Aggregate,
							"key1",
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				projectID:    "project1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ro-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.BulkRemoveProjectMetadata(tt.args.ctx, tt.args.projectID, "ro-1", tt.args.metadataList...)
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
