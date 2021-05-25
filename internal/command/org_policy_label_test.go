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
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/static/mock"
)

func TestCommandSide_AddLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:    "",
					BackgroundColor: "#ffffff",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "labelpolicy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
				ctx:   context.Background(),
				orgID: "org1",
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
								org.NewLabelPolicyAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
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
				ctx:   context.Background(),
				orgID: "org1",
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
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
			got, err := r.AddLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangeLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LabelPolicy{
					PrimaryColor:    "#ffffff",
					BackgroundColor: "#ffffff",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "labelpolicy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								newLabelPolicyChangedEvent(
									context.Background(),
									"org1",
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
				ctx:   context.Background(),
				orgID: "org1",
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
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
			got, err := r.ChangeLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ActivateLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org id missing, invalid argument error",
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
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "activate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyActivatedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			_, err := r.ActivateLabelPolicy(tt.args.ctx, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "org id missing, invalid argument error",
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
			name: "labelpolicy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			_, err := r.RemoveLabelPolicy(tt.args.ctx, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_AddLogoLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "storage key empty, invalid argument error",
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
				orgID:      "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyLogoAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
				static:     tt.fields.storage,
			}
			got, err := r.AddLogoLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.storageKey)
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

func TestCommandSide_RemoveLogoLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
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
				orgID:      "org1",
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLabelPolicyLogoAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLabelPolicyLogoRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
				static:     tt.fields.storage,
			}
			got, err := r.RemoveLogoLabelPolicy(tt.args.ctx, tt.args.orgID)
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

func TestCommandSide_AddIconLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "storage key empty, invalid argument error",
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
				orgID:      "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyIconAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
			got, err := r.AddIconLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.storageKey)
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

func TestCommandSide_RemoveIconLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "",
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
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon added, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLabelPolicyIconAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLabelPolicyIconRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
				static:     tt.fields.storage,
			}
			got, err := r.RemoveIconLabelPolicy(tt.args.ctx, tt.args.orgID)
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

func TestCommandSide_AddLogoDarkLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "storage key empty, invalid argument error",
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
				orgID:      "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
			got, err := r.AddLogoDarkLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.storageKey)
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

func TestCommandSide_RemoveLogoDarkLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
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
				orgID:      "org1",
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsNotFound,
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLabelPolicyLogoDarkRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
				static:     tt.fields.storage,
			}
			got, err := r.RemoveLogoDarkLabelPolicy(tt.args.ctx, tt.args.orgID)
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

func TestCommandSide_AddIconDarkLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "storage key empty, invalid argument error",
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
				orgID:      "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyIconDarkAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
			got, err := r.AddIconDarkLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.storageKey)
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

func TestCommandSide_RemoveIconDarkLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "orgID empty, invalid argument error",
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
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "icon dark added, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLabelPolicyIconDarkAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLabelPolicyIconDarkRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
				static:     tt.fields.storage,
			}
			got, err := r.RemoveIconDarkLabelPolicy(tt.args.ctx, tt.args.orgID)
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

func TestCommandSide_AddFontLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		orgID      string
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
			name: "orgID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:        context.Background(),
				storageKey: "key",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "storage key empty, invalid argument error",
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
				orgID:      "org1",
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								org.NewLabelPolicyFontAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				orgID:      "org1",
				storageKey: "key",
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
			got, err := r.AddFontLabelPolicy(tt.args.ctx, tt.args.orgID, tt.args.storageKey)
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

func TestCommandSide_RemoveFontLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "orgID empty, invalid argument error",
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
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
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
							org.NewLabelPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
							org.NewLabelPolicyFontAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLabelPolicyFontRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"key",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
				static:     tt.fields.storage,
			}
			got, err := r.RemoveFontLabelPolicy(tt.args.ctx, tt.args.orgID)
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

func newLabelPolicyChangedEvent(ctx context.Context, orgID, primaryColor, backgroundColor, warnColor, fontColor, primaryColorDark, backgroundColorDark, warnColorDark, fontColorDark string, hideLoginNameSuffix, errMsgPopup, disableWatermark bool) *org.LabelPolicyChangedEvent {
	event, _ := org.NewLabelPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
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
