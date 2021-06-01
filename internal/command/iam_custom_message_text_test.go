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
	"github.com/caos/zitadel/internal/repository/iam"
)

func TestCommandSide_SetDefaultMessageText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		config *domain.CustomMessageText
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
			name: "invalid custom text, error",
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
			name: "custom text set all fields, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageGreeting,
									"Greeting",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageSubject,
									"Subject",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageTitle,
									"Title",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessagePreHeader,
									"PreHeader",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageText,
									"Text",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageButtonText,
									"ButtonText",
									language.English,
								),
							),
							eventFromEventPusher(
								iam.NewCustomTextSetEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"Template",
									domain.MessageFooterText,
									"FooterText",
									language.English,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.CustomMessageText{
					MessageTextType: "Template",
					Language:        language.English,
					Greeting:        "Greeting",
					Subject:         "Subject",
					Title:           "Title",
					PreHeader:       "PreHeader",
					Text:            "Text",
					ButtonText:      "ButtonText",
					FooterText:      "FooterText",
				},
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
			got, err := r.SetDefaultMessageText(tt.args.ctx, tt.args.config)
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
