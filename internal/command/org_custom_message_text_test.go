package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_SetCustomMessageText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		config        *domain.CustomMessageText
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
			name: "no resource owner, error",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx:    context.Background(),
				config: &domain.CustomMessageText{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty message type, error",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomMessageText{
					Language: AllowedLanguage,
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "empty custom message text, success",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomMessageText{
					MessageTextType: "Some type", // TODO: check the type!
					Language:        AllowedLanguage,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "undefined language, error",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config:        &domain.CustomMessageText{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "unsupported language, error",
			fields: fields{
				eventstore: eventstoreExpect(t),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomMessageText{
					Language: UnsupportedLanguage,
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "custom text set all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageGreeting,
							"Greeting",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageSubject,
							"Subject",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageTitle,
							"Title",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessagePreHeader,
							"PreHeader",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageText,
							"Text",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageButtonText,
							"ButtonText",
							language.English,
						),

						org.NewCustomTextSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageFooterText,
							"Footer",
							language.English,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomMessageText{
					MessageTextType: "Template",
					Language:        language.English,
					Greeting:        "Greeting",
					Subject:         "Subject",
					Title:           "Title",
					PreHeader:       "PreHeader",
					Text:            "Text",
					ButtonText:      "ButtonText",
					FooterText:      "Footer",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "custom text remove all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageGreeting,
								"Greeting",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageSubject,
								"Subject",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageTitle,
								"Title",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessagePreHeader,
								"PreHeader",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageText,
								"Text",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageButtonText,
								"ButtonText",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageFooterText,
								"Footer",
								language.English,
							),
						),
					),
					expectPush(
						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageGreeting,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageSubject,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageTitle,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessagePreHeader,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageText,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageButtonText,
							language.English,
						),

						org.NewCustomTextRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							domain.MessageFooterText,
							language.English,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CustomMessageText{
					MessageTextType: "Template",
					Language:        language.English,
					Greeting:        "",
					Subject:         "",
					Title:           "",
					PreHeader:       "",
					Text:            "",
					ButtonText:      "",
					FooterText:      "",
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
			got, err := r.SetOrgMessageText(tt.args.ctx, tt.args.resourceOwner, tt.args.config)
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

func TestCommandSide_RemoveCustomMessageText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		mailTextType  string
		lang          language.Tag
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
			name: "no resource owner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:          context.Background(),
				mailTextType: "Template",
				lang:         language.English,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no mail text type owner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				lang:          language.English,
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "no mail text type owner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				mailTextType:  "Template",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "custom text remove all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageGreeting,
								"Greeting",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageSubject,
								"Subject",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageTitle,
								"Title",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessagePreHeader,
								"PreHeader",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageText,
								"Text",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageButtonText,
								"ButtonText",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageFooterText,
								"Footer",
								language.English,
							),
						),
					),
					expectPush(
						org.NewCustomTextTemplateRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							language.English,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				mailTextType:  "Template",
				lang:          language.English,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove unsupported language ok, especially because we never validated whether a language is supported in previous ZITADEL versions",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"Template",
								domain.MessageGreeting,
								"Greeting",
								UnsupportedLanguage,
							),
						),
					),
					expectPush(
						org.NewCustomTextTemplateRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"Template",
							UnsupportedLanguage,
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				mailTextType:  "Template",
				lang:          UnsupportedLanguage,
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
			_, err := r.RemoveOrgMessageTexts(tt.args.ctx, tt.args.resourceOwner, tt.args.mailTextType, tt.args.lang)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
