package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetUserMetadata(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx      context.Context
			orgID    string
			userID   string
			metadata *domain.Metadata
			check    PermissionCheck
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
			name: "user not existing, pre condition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
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
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key: "key",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
				check: func(resourceOwner, aggregateID string) error {
					return zerrors.ThrowPermissionDenied(nil, "id", "permission denied")
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(),
					expectPush(
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
							[]byte("value"),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: []byte("value"),
					State: domain.MetadataStateActive,
				},
			},
		},
		{
			name: "add metadata, reset, invalid",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							user.NewMetadataRemovedAllEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key: "key",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata with same key, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
					expectPush(
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
							[]byte("value2"),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value2"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: []byte("value2"),
					State: domain.MetadataStateActive,
				},
			},
		},
		{
			name: "add metadata with same key and value, ok (ignore)",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: []byte("value"),
					State: domain.MetadataStateActive,
				},
			},
		},
		{
			name: "add deleted metadata with same value, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							user.NewMetadataRemovedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
							[]byte("value"),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadata: &domain.Metadata{
					Key:   "key",
					Value: []byte("value"),
				},
			},
			res: res{
				want: &domain.Metadata{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.SetUserMetadata(tt.args.ctx, tt.args.metadata, tt.args.userID, tt.args.orgID, tt.args.check)
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

func TestCommandSide_BulkSetUserMetadata(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			userID       string
			check        PermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user not existing, pre condition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
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
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
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
			name: "add metadata, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				check: func(resourceOwner, aggregateID string) error {
					return zerrors.ThrowPermissionDenied(nil, "id", "permission-denied")
				},
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(),
					expectPush(
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
							[]byte("value"),
						),
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key1",
							[]byte("value1"),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
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
		{
			name: "re add deleted metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							)),
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key1",
								[]byte("value1"),
							)),
						eventFromEventPusher(
							user.NewMetadataRemovedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
							)),
					),
					expectPush(
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
							[]byte("value"),
						),
						user.NewMetadataSetEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key2",
							[]byte("value2"),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metadataList: []*domain.Metadata{
					{Key: "key", Value: []byte("value")},
					{Key: "key1", Value: []byte("value1")},
					{Key: "key2", Value: []byte("value2")},
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.BulkSetUserMetadata(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.check, tt.args.metadataList...)
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

func TestCommandSide_UserRemoveMetadata(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx         context.Context
			orgID       string
			userID      string
			metadataKey string
			check       PermissionCheck
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
			name: "user not existing, pre condition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metadataKey: "key",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metadataKey: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "meta data not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metadataKey: "key",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove metadata, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metadataKey: "key",
				check: func(resourceOwner, aggregateID string) error {
					return zerrors.ThrowPermissionDenied(nil, "id", "permission denied")
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
					),
					expectPush(
						user.NewMetadataRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.RemoveUserMetadata(tt.args.ctx, tt.args.metadataKey, tt.args.userID, tt.args.orgID, tt.args.check)
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

func TestCommandSide_BulkRemoveUserMetadata(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			userID       string
			check        PermissionCheck
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "user not existing, pre condition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "remove metadata keys not existing, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
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
				userID:       "user1",
				metadataList: []string{"key", "key1"},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
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
				userID:       "user1",
				metadataList: []string{""},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "remove metadata, no permission",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
				metadataList: []string{"key", "key1"},
				check: func(resourceOwner, aggregateID string) error {
					return zerrors.ThrowPermissionDenied(nil, "id", "permission denied")
				},
			},
			res: res{
				err: zerrors.IsPermissionDenied,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.Und,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								[]byte("value"),
							),
						),
						eventFromEventPusher(
							user.NewMetadataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key1",
								[]byte("value1"),
							),
						),
					),
					expectPush(
						user.NewMetadataRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key",
						),
						user.NewMetadataRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"key1",
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
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
				eventstore: tt.fields.eventstore(t),
			}
			got, err := r.BulkRemoveUserMetadata(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.check, tt.args.metadataList...)
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
