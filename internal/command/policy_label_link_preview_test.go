package command

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/static/mock"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeDefaultLabelPolicy_LinkPreview(t *testing.T) {
	type res struct {
		want *domain.LabelPolicy
		err  func(error) bool
	}
	tests := []struct {
		name        string
		eventstore  func(*testing.T) *eventstore.Eventstore
		policy      *domain.LabelPolicy
		linkPreview LabelPolicyLinkPreviewUpdate
		res         res
	}{
		{
			name:        "title too long",
			eventstore:  expectEventstore(),
			policy:      whiteLabelPolicy(),
			linkPreview: LabelPolicyLinkPreviewUpdate{Title: gu.Ptr(strings.Repeat("a", 121))},
			res:         res{err: zerrors.IsErrorInvalidArgument},
		},
		{
			name:        "description too long",
			eventstore:  expectEventstore(),
			policy:      whiteLabelPolicy(),
			linkPreview: LabelPolicyLinkPreviewUpdate{Description: gu.Ptr(strings.Repeat("a", 301))},
			res:         res{err: zerrors.IsErrorInvalidArgument},
		},
		{
			name: "link preview text changed, ok",
			eventstore: expectEventstore(
				expectFilter(eventFromEventPusher(defaultLabelPolicyAddedEvent("", ""))),
				expectPush(newLinkPreviewChangedEvent(context.Background(), "Acme", "Sign in to Acme")),
			),
			policy:      whiteLabelPolicy(),
			linkPreview: LabelPolicyLinkPreviewUpdate{Title: gu.Ptr("Acme"), Description: gu.Ptr("Sign in to Acme")},
			res:         res{want: wantWhiteLabelPolicy("Acme", "Sign in to Acme")},
		},
		{
			name: "missing link preview text keeps the current values",
			eventstore: expectEventstore(
				expectFilter(eventFromEventPusher(defaultLabelPolicyAddedEvent("Acme", "Sign in to Acme"))),
				expectPush(newDefaultLabelPolicyChangedEventWith(policy.ChangePrimaryColor("#000000"))),
			),
			policy: func() *domain.LabelPolicy {
				p := whiteLabelPolicy()
				p.PrimaryColor = "#000000"
				return p
			}(),
			res: res{want: func() *domain.LabelPolicy {
				p := wantWhiteLabelPolicy("Acme", "Sign in to Acme")
				p.PrimaryColor = "#000000"
				return p
			}()},
		},
		{
			name: "empty link preview text clears the values",
			eventstore: expectEventstore(
				expectFilter(eventFromEventPusher(defaultLabelPolicyAddedEvent("Acme", "Sign in to Acme"))),
				expectPush(newLinkPreviewChangedEvent(context.Background(), "", "")),
			),
			policy:      whiteLabelPolicy(),
			linkPreview: LabelPolicyLinkPreviewUpdate{Title: gu.Ptr(""), Description: gu.Ptr("")},
			res:         res{want: wantWhiteLabelPolicy("", "")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{eventstore: tt.eventstore(t)}
			got, err := r.ChangeDefaultLabelPolicy(context.Background(), tt.policy, tt.linkPreview)
			if tt.res.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.res.want, got)
				return
			}
			if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_OrgLabelPolicy_LinkPreview(t *testing.T) {
	orgAdded := func(title, description string) *org.LabelPolicyAddedEvent {
		return org.NewLabelPolicyAddedEvent(context.Background(),
			&org.NewAggregate("org1").Aggregate,
			"#ffffff", "#ffffff", "#ffffff", "#ffffff",
			"#ffffff", "#ffffff", "#ffffff", "#ffffff",
			true, true, true,
			domain.LabelPolicyThemeAuto, title, description,
		)
	}

	t.Run("add stores the link preview text", func(t *testing.T) {
		r := &Commands{eventstore: expectEventstore(
			expectFilter(),
			expectPush(orgAdded("Acme", "Sign in to Acme")),
		)(t)}
		p := whiteLabelPolicy()
		p.LinkPreviewTitle = "Acme"
		p.LinkPreviewDescription = "Sign in to Acme"
		got, err := r.AddLabelPolicy(context.Background(), "org1", p)
		assert.NoError(t, err)
		assert.Equal(t, "Acme", got.LinkPreviewTitle)
		assert.Equal(t, "Sign in to Acme", got.LinkPreviewDescription)
	})

	t.Run("change without link preview text keeps the current values", func(t *testing.T) {
		changed, _ := org.NewLabelPolicyChangedEvent(context.Background(),
			&org.NewAggregate("org1").Aggregate,
			[]policy.LabelPolicyChanges{policy.ChangePrimaryColor("#000000")},
		)
		r := &Commands{eventstore: expectEventstore(
			expectFilter(eventFromEventPusher(orgAdded("Acme", "Sign in to Acme"))),
			expectPush(changed),
		)(t)}
		p := whiteLabelPolicy()
		p.PrimaryColor = "#000000"
		got, err := r.ChangeLabelPolicy(context.Background(), "org1", p, LabelPolicyLinkPreviewUpdate{})
		assert.NoError(t, err)
		assert.Equal(t, "Acme", got.LinkPreviewTitle)
		assert.Equal(t, "Sign in to Acme", got.LinkPreviewDescription)
	})

	t.Run("change updates the link preview text", func(t *testing.T) {
		changed, _ := org.NewLabelPolicyChangedEvent(context.Background(),
			&org.NewAggregate("org1").Aggregate,
			[]policy.LabelPolicyChanges{policy.ChangeLinkPreviewTitle("Roadrunner")},
		)
		r := &Commands{eventstore: expectEventstore(
			expectFilter(eventFromEventPusher(orgAdded("Acme", "Sign in to Acme"))),
			expectPush(changed),
		)(t)}
		got, err := r.ChangeLabelPolicy(context.Background(), "org1", whiteLabelPolicy(), LabelPolicyLinkPreviewUpdate{Title: gu.Ptr("Roadrunner")})
		assert.NoError(t, err)
		assert.Equal(t, "Roadrunner", got.LinkPreviewTitle)
		assert.Equal(t, "Sign in to Acme", got.LinkPreviewDescription)
	})
}

func whiteLabelPolicy() *domain.LabelPolicy {
	return &domain.LabelPolicy{
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
	}
}

func wantWhiteLabelPolicy(title, description string) *domain.LabelPolicy {
	p := whiteLabelPolicy()
	p.ObjectRoot = models.ObjectRoot{
		AggregateID:   "INSTANCE",
		ResourceOwner: "INSTANCE",
		InstanceID:    "INSTANCE",
	}
	p.LinkPreviewTitle = title
	p.LinkPreviewDescription = description
	return p
}

func newDefaultLabelPolicyChangedEventWith(changes ...policy.LabelPolicyChanges) *instance.LabelPolicyChangedEvent {
	event, _ := instance.NewLabelPolicyChangedEvent(context.Background(), &instance.NewAggregate("INSTANCE").Aggregate, changes)
	return event
}

func newLinkPreviewChangedEvent(ctx context.Context, title, description string) *instance.LabelPolicyChangedEvent {
	event, _ := instance.NewLabelPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.LabelPolicyChanges{
			policy.ChangeLinkPreviewTitle(title),
			policy.ChangeLinkPreviewDescription(description),
		},
	)
	return event
}

func TestCommandSide_AddLinkPreviewImageDefaultLabelPolicy(t *testing.T) {
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
				eventstore: eventstoreExpect(t, expectFilter()),
			},
			args: args{
				ctx:    context.Background(),
				upload: newPreviewUpload("IAM"),
			},
			res: res{err: zerrors.IsNotFound},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(eventFromEventPusher(defaultLabelPolicyAddedEvent("", ""))),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx:    context.Background(),
				upload: newPreviewUpload("IAM"),
			},
			res: res{err: zerrors.IsInternal},
		},
		{
			name: "image added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(eventFromEventPusher(defaultLabelPolicyAddedEvent("", ""))),
					expectPush(
						instance.NewLabelPolicyLinkPreviewImageAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"preview.png",
						),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx:    context.Background(),
				upload: newPreviewUpload("IAM"),
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "INSTANCE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.AddLinkPreviewImageDefaultLabelPolicy(tt.args.ctx, tt.args.upload)
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

func TestCommandSide_RemoveLinkPreviewImageDefaultLabelPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "label policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(t, expectFilter()),
			},
			res: res{err: zerrors.IsNotFound},
		},
		{
			name: "image removed, ok",
			fields: fields{
				storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(defaultLabelPolicyAddedEvent("", "")),
						eventFromEventPusher(
							instance.NewLabelPolicyLinkPreviewImageAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"key",
							),
						),
					),
					expectPush(
						instance.NewLabelPolicyLinkPreviewImageRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"key",
						),
					),
				),
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "INSTANCE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.RemoveLinkPreviewImageDefaultLabelPolicy(context.Background())
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

func TestCommandSide_AddLinkPreviewImageLabelPolicy(t *testing.T) {
	tests := []struct {
		name       string
		eventstore *eventstore.Eventstore
		storage    static.Storage
		orgID      string
		upload     *AssetUpload
		want       *domain.ObjectDetails
		err        func(error) bool
	}{
		{
			name:       "orgID empty, invalid argument error",
			eventstore: eventstoreExpect(t),
			err:        zerrors.IsErrorInvalidArgument,
		},
		{
			name:       "label policy not existing, not found error",
			eventstore: eventstoreExpect(t, expectFilter()),
			orgID:      "org1",
			upload:     newPreviewUpload("org1"),
			err:        zerrors.IsNotFound,
		},
		{
			name: "image added, ok",
			eventstore: eventstoreExpect(
				t,
				expectFilter(eventFromEventPusher(orgLabelPolicyAddedEvent())),
				expectPush(
					org.NewLabelPolicyLinkPreviewImageAddedEvent(context.Background(),
						&org.NewAggregate("org1").Aggregate,
						"preview.png",
					),
				),
			),
			storage: mock.NewStorage(t).ExpectPutObject(),
			orgID:   "org1",
			upload:  newPreviewUpload("org1"),
			want:    &domain.ObjectDetails{ResourceOwner: "org1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{eventstore: tt.eventstore, static: tt.storage}
			got, err := r.AddLinkPreviewImageLabelPolicy(context.Background(), tt.orgID, tt.upload)
			if tt.err == nil {
				assert.NoError(t, err)
				assertObjectDetails(t, tt.want, got)
				return
			}
			if !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCommandSide_RemoveLinkPreviewImageLabelPolicy(t *testing.T) {
	tests := []struct {
		name       string
		eventstore *eventstore.Eventstore
		storage    static.Storage
		orgID      string
		want       *domain.ObjectDetails
		err        func(error) bool
	}{
		{
			name:       "orgID empty, invalid argument error",
			eventstore: eventstoreExpect(t),
			err:        zerrors.IsErrorInvalidArgument,
		},
		{
			name:    "image removed, ok",
			storage: mock.NewMockStorage(gomock.NewController(t)).ExpectRemoveObjectNoError(),
			eventstore: eventstoreExpect(
				t,
				expectFilter(
					eventFromEventPusher(orgLabelPolicyAddedEvent()),
					eventFromEventPusher(
						org.NewLabelPolicyLinkPreviewImageAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"key",
						),
					),
				),
				expectPush(
					org.NewLabelPolicyLinkPreviewImageRemovedEvent(context.Background(),
						&org.NewAggregate("org1").Aggregate,
						"key",
					),
				),
			),
			orgID: "org1",
			want:  &domain.ObjectDetails{ResourceOwner: "org1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{eventstore: tt.eventstore, static: tt.storage}
			got, err := r.RemoveLinkPreviewImageLabelPolicy(context.Background(), tt.orgID)
			if tt.err == nil {
				assert.NoError(t, err)
				assertObjectDetails(t, tt.want, got)
				return
			}
			if !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newPreviewUpload(resourceOwner string) *AssetUpload {
	return &AssetUpload{
		ResourceOwner: resourceOwner,
		ObjectName:    "preview.png",
		ContentType:   "image",
		ObjectType:    static.ObjectTypeStyling,
		File:          bytes.NewReader([]byte("test")),
		Size:          4,
	}
}

func defaultLabelPolicyAddedEvent(linkPreviewTitle, linkPreviewDescription string) *instance.LabelPolicyAddedEvent {
	return instance.NewLabelPolicyAddedEvent(context.Background(),
		&instance.NewAggregate("INSTANCE").Aggregate,
		"#ffffff", "#ffffff", "#ffffff", "#ffffff",
		"#ffffff", "#ffffff", "#ffffff", "#ffffff",
		true, true, true,
		domain.LabelPolicyThemeAuto, linkPreviewTitle, linkPreviewDescription,
	)
}

func orgLabelPolicyAddedEvent() *org.LabelPolicyAddedEvent {
	return org.NewLabelPolicyAddedEvent(context.Background(),
		&org.NewAggregate("org1").Aggregate,
		"#ffffff", "#ffffff", "#ffffff", "#ffffff",
		"#ffffff", "#ffffff", "#ffffff", "#ffffff",
		true, true, true,
		domain.LabelPolicyThemeAuto, "", "",
	)
}
