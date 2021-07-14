package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_SetMetaData(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx      context.Context
			orgID    string
			userID   string
			metaData *domain.MetaData
		}
	)
	type res struct {
		want *domain.MetaData
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metaData: &domain.MetaData{
					Key:   "key",
					Value: "value",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				metaData: &domain.MetaData{
					Key: "key",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMetaDataSetEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key",
									"value",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metaData: &domain.MetaData{
					Key:   "key",
					Value: "value",
				},
			},
			res: res{
				want: &domain.MetaData{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Key:   "key",
					Value: "value",
					State: domain.MetaDataStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetUserMetaData(tt.args.ctx, tt.args.metaData, tt.args.userID, tt.args.orgID)
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

func TestCommandSide_BulkSetMetaData(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			userID       string
			metaDataList []*domain.MetaData
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
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "user not existing, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metaDataList: []*domain.MetaData{
					{Key: "key", Value: "value"},
					{Key: "key1", Value: "value1"},
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				metaDataList: []*domain.MetaData{
					{Key: "key"},
					{Key: "key1"},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMetaDataSetEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key",
									"value",
								),
							),
							eventFromEventPusher(
								user.NewMetaDataSetEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key1",
									"value1",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
				metaDataList: []*domain.MetaData{
					{Key: "key", Value: "value"},
					{Key: "key1", Value: "value1"},
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
			got, err := r.BulkSetUserMetaData(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.metaDataList...)
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

func TestCommandSide_UserRemoveMetaData(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx         context.Context
			orgID       string
			userID      string
			metaDataKey string
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
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metaDataKey: "key",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
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
				userID:      "user1",
				metaDataKey: "",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "meta data not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				metaDataKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								"value",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMetaDataRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:         context.Background(),
				orgID:       "org1",
				userID:      "user1",
				metaDataKey: "key",
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
			got, err := r.RemoveUserMetaData(tt.args.ctx, tt.args.metaDataKey, tt.args.userID, tt.args.orgID)
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

func TestCommandSide_BulkRemoveMetaData(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx          context.Context
			orgID        string
			userID       string
			metaDataList []string
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
		//{
		//	name: "empty meta data list, pre condition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//		),
		//	},
		//	args: args{
		//		ctx:    context.Background(),
		//		orgID:  "org1",
		//		userID: "user1",
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		//{
		//	name: "user not existing, pre condition error",
		//	fields: fields{
		//		eventstore: eventstoreExpect(
		//			t,
		//			expectFilter(),
		//		),
		//	},
		//	args: args{
		//		ctx:    context.Background(),
		//		orgID:  "org1",
		//		userID: "user1",
		//		metaDataList: []string{"key", "key1"},
		//	},
		//	res: res{
		//		err: caos_errs.IsPreconditionFailed,
		//	},
		//},
		{
			name: "remove metadata keys not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								"value",
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
				metaDataList: []string{"key", "key1"},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "invalid metadata, pre condition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								"value",
							),
						),
						eventFromEventPusher(
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key1",
								"value1",
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
				metaDataList: []string{""},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "remove metadata, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key",
								"value",
							),
						),
						eventFromEventPusher(
							user.NewMetaDataSetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"key1",
								"value1",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								user.NewMetaDataRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key",
								),
							),
							eventFromEventPusher(
								user.NewMetaDataRemovedEvent(context.Background(),
									&user.NewAggregate("user1", "org1").Aggregate,
									"key1",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				userID:       "user1",
				metaDataList: []string{"key", "key1"},
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
			got, err := r.BulkRemoveUserMetaData(tt.args.ctx, tt.args.userID, tt.args.orgID, tt.args.metaDataList...)
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
