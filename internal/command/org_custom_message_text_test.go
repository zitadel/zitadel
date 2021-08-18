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
	"github.com/caos/zitadel/internal/repository/org"
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
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				config: &domain.CustomMessageText{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid custom text, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config:        &domain.CustomMessageText{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "custom text set all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageGreeting,
									"Greeting",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageSubject,
									"Subject",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageTitle,
									"Title",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessagePreHeader,
									"PreHeader",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageText,
									"Text",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageButtonText,
									"ButtonText",
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextSetEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageFooterText,
									"Footer",
									language.English,
								),
							),
						},
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
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageGreeting,
								"Greeting",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageSubject,
								"Subject",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageTitle,
								"Title",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessagePreHeader,
								"PreHeader",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageText,
								"Text",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageButtonText,
								"ButtonText",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageFooterText,
								"Footer",
								language.English,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageGreeting,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageSubject,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageTitle,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessagePreHeader,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageText,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageButtonText,
									language.English,
								),
							),
							eventFromEventPusher(
								org.NewCustomTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									domain.MessageFooterText,
									language.English,
								),
							),
						},
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
				err: caos_errs.IsErrorInvalidArgument,
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
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageGreeting,
								"Greeting",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageSubject,
								"Subject",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageTitle,
								"Title",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessagePreHeader,
								"PreHeader",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageText,
								"Text",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageButtonText,
								"ButtonText",
								language.English,
							),
						),
						eventFromEventPusher(
							org.NewCustomTextSetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"Template",
								domain.MessageFooterText,
								"Footer",
								language.English,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewCustomTextTemplateRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"Template",
									language.English,
								),
							),
						},
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
