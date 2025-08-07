package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/static/mock"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
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
					ThemeMode:           domain.LabelPolicyThemeAuto,
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
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
							false,
							domain.LabelPolicyThemeDark),
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
					ThemeMode:           domain.LabelPolicyThemeDark,
				},
			},
			res: res{
				want: &domain.LabelPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
						InstanceID:    "INSTANCE",
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
					ThemeMode:           domain.LabelPolicyThemeDark,
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
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "activated, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyActivatedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
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
		ctx    context.Context
		upload *AssetUpload
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
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "text/css",
					ObjectType:    static.ObjectTypeStyling,
					File:          nil,
					Size:          0,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "logo added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyLogoAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"logo",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			got, err := r.AddLogoDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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
				err: zerrors.IsNotFound,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyLogoAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				err: zerrors.IsInternal,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyLogoAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyLogoRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddIconDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx    context.Context
		upload *AssetUpload
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
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "icon added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyIconAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"icon",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			got, err := r.AddIconDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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
				err: zerrors.IsNotFound,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyIconAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyIconRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddLogoDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx        context.Context
		instanceID string
		upload     *AssetUpload
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
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "logo dark added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"logo",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			got, err := r.AddLogoDarkDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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
				err: zerrors.IsNotFound,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				err: zerrors.IsInternal,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyLogoDarkAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyLogoDarkRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddIconDarkDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx    context.Context
		upload *AssetUpload
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
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "icon dark added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyIconDarkAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"icon",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "icon",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			got, err := r.AddIconDarkDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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
				err: zerrors.IsNotFound,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyIconDarkAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				err: zerrors.IsInternal,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyIconDarkAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyIconDarkRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddFontDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx    context.Context
		upload *AssetUpload
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
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "font",
					ContentType:   "ttf",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "font",
					ContentType:   "ttf",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: zerrors.IsInternal,
			},
		},
		{
			name: "font added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyFontAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"font",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx: context.Background(),
				upload: &AssetUpload{
					ResourceOwner: "IAM",
					ObjectName:    "font",
					ContentType:   "ttf",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
			got, err := r.AddFontDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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
				err: zerrors.IsNotFound,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyFontAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
				err: zerrors.IsInternal,
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
							instance.NewLabelPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
								domain.LabelPolicyThemeAuto,
							),
						),
						eventFromEventPusher(
							instance.NewLabelPolicyFontAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyFontRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func newDefaultLabelPolicyChangedEvent(ctx context.Context, primaryColor, backgroundColor, warnColor, fontColor, primaryColorDark, backgroundColorDark, warnColorDark, fontColorDark string, hideLoginNameSuffix, errMsgPopup, disableWatermark bool, theme domain.LabelPolicyThemeMode) *instance.LabelPolicyChangedEvent {
	event, _ := instance.NewLabelPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
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
			policy.ChangeThemeMode(theme),
		},
	)
	return event
}
