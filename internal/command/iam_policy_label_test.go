package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/static/mock"
)

func TestCommandSide_AddDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LabelPolicy
	}
	type res struct {
		want *domain.LabelPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "labelpolicy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:        "#ffffff",
					BackgroundColor:     "#ffffff",
					WarnColor:           "#ffffff",
					FontColor:           "#ffffff",
					PrimaryColorDark:    "#ffffff",
					BackgroundColorDark: "#ffffff",
					WarnColorDark:       "#ffffff",
					FontColorDark:       "#ffffff",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"#ffffff",
									"#ffffff",
									"#ffffff",
									"#ffffff",
									"#ffffff",
									"#ffffff",
									"#ffffff",
									"#ffffff",
									true,
									true,
									true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:        "#ffffff",
					BackgroundColor:     "#ffffff",
					WarnColor:           "#ffffff",
					FontColor:           "#ffffff",
					PrimaryColorDark:    "#ffffff",
					BackgroundColorDark: "#ffffff",
					WarnColorDark:       "#ffffff",
					FontColorDark:       "#ffffff",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
				},
			},
			res: res{
				want: &domain.LabelPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					PrimaryColor:        "#ffffff",
					BackgroundColor:     "#ffffff",
					WarnColor:           "#ffffff",
					FontColor:           "#ffffff",
					PrimaryColorDark:    "#ffffff",
					BackgroundColorDark: "#ffffff",
					WarnColorDark:       "#ffffff",
					FontColorDark:       "#ffffff",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultLabelPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LabelPolicy
	}
	type res struct {
		want *domain.LabelPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "labelpolicy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:    "#ffffff",
					BackgroundColor: "#ffffff",
					WarnColor:       "#ffffff",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:        "#ffffff",
					BackgroundColor:     "#ffffff",
					WarnColor:           "#ffffff",
					FontColor:           "#ffffff",
					PrimaryColorDark:    "#ffffff",
					BackgroundColorDark: "#ffffff",
					WarnColorDark:       "#ffffff",
					FontColorDark:       "#ffffff",
					HideLoginNameSuffix: true,
					ErrorMsgPopup:       true,
					DisableWatermark:    true,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultLabelPolicyChangedEvent(
									context.Background(),
									"#000000",
									"#000000",
									"#000000",
									"#000000",
									"#000000",
									"#000000",
									"#000000",
									"#000000",
									false,
									false,
									false),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:        "#000000",
					BackgroundColor:     "#000000",
					WarnColor:           "#000000",
					FontColor:           "#000000",
					PrimaryColorDark:    "#000000",
					BackgroundColorDark: "#000000",
					WarnColorDark:       "#000000",
					FontColorDark:       "#000000",
					HideLoginNameSuffix: false,
					ErrorMsgPopup:       false,
					DisableWatermark:    false,
				},
			},
			res: res{
				want: &domain.LabelPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					PrimaryColor:        "#000000",
					BackgroundColor:     "#000000",
					WarnColor:           "#000000",
					FontColor:           "#000000",
					PrimaryColorDark:    "#000000",
					BackgroundColorDark: "#000000",
					WarnColorDark:       "#000000",
					FontColorDark:       "#000000",
					HideLoginNameSuffix: false,
					ErrorMsgPopup:       false,
					DisableWatermark:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultLabelPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ActivateDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "activated, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyActivatedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ActivateDefaultLabelPolicy(tt.args.ctx)
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

func TestCommandSide_AddLogoDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		storageKey string
	}
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
			name: "storage key empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "logo added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyLogoAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddLogoDefaultLabelPolicy(tt.args.ctx, tt.args.storageKey)
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

func TestCommandSide_RemoveLogoDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},

		{
			name: "asset remove error, internal error",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyLogoAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "logo added, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyLogoAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyLogoRemovedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveLogoDefaultLabelPolicy(tt.args.ctx)
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

func TestCommandSide_AddIconDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		storageKey string
	}
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
			name: "storage key empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyIconAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddIconDefaultLabelPolicy(tt.args.ctx, tt.args.storageKey)
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

func TestCommandSide_RemoveIconDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon removed, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyIconAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyIconRemovedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveIconDefaultLabelPolicy(tt.args.ctx)
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

func TestCommandSide_AddLogoDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		storageKey string
	}
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
			name: "storage key empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "logo dark added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddLogoDarkDefaultLabelPolicy(tt.args.ctx, tt.args.storageKey)
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

func TestCommandSide_RemoveLogoDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "logo dark removed, not ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "logo dark removed, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyLogoDarkRemovedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveLogoDarkDefaultLabelPolicy(tt.args.ctx)
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

func TestCommandSide_AddIconDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		storageKey string
	}
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
			name: "storage key empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon dark added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyIconDarkAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddIconDarkDefaultLabelPolicy(tt.args.ctx, tt.args.storageKey)
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

func TestCommandSide_RemoveIconDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon dark removed, not ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyIconDarkAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "icon dark removed, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyIconDarkAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyIconDarkRemovedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveIconDarkDefaultLabelPolicy(tt.args.ctx)
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

func TestCommandSide_AddFontDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		storageKey string
	}
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
			name: "storage key empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "font added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyFontAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddFontDefaultLabelPolicy(tt.args.ctx, tt.args.storageKey)
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

func TestCommandSide_RemoveFontDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx context.Context
	}
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
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "font remove from storage not possible, internla error",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyFontAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "font added, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
							),
						),
						eventFromEventPusher(
							iam.NewLabelPolicyFontAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewLabelPolicyFontRemovedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveFontDefaultLabelPolicy(tt.args.ctx)
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

func newDefaultLabelPolicyChangedEvent(ctx context.Context, primaryColor, backgroundColor, warnColor, fontColor, primaryColorDark, backgroundColorDark, warnColorDark, fontColorDark string, hideLoginNameSuffix, errMsgPopup, disableWatermark bool) *iam.LabelPolicyChangedEvent {
	event, _ := iam.NewLabelPolicyChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		[]policy.LabelPolicyChanges{
			policy.ChangePrimaryColor(primaryColor),
			policy.ChangeBackgroundColor(backgroundColor),
			policy.ChangeWarnColor(warnColor),
			policy.ChangeFontColor(fontColor),
			policy.ChangePrimaryColorDark(primaryColorDark),
			policy.ChangeBackgroundColorDark(backgroundColorDark),
			policy.ChangeWarnColorDark(warnColorDark),
			policy.ChangeFontColorDark(fontColorDark),
			policy.ChangeHideLoginNameSuffix(hideLoginNameSuffix),
			policy.ChangeErrorMsgPopup(errMsgPopup),
			policy.ChangeDisableWatermark(disableWatermark),
		},
	)
	return event
}
