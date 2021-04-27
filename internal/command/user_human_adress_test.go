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

func TestCommandSide_ChangeHumanAddress(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		address       *domain.Address
		resourceOwner string
	}
	type res struct {
		want *domain.Address
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "user not existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Address{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					Country:    "Switzerland",
					Locality:   "St. Gallen",
					PostalCode: "9000",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "address not changed, precondition error",
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
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
						eventFromEventPusher(
							newAddressChangedEvent(context.Background(),
								"user1", "org1",
								"country",
								"locality",
								"postalcode",
								"region",
								"street",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Address{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					Country:       "country",
					Locality:      "locality",
					PostalCode:    "postalcode",
					Region:        "region",
					StreetAddress: "street",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "address changed, ok",
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
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newAddressChangedEvent(context.Background(),
									"user1", "org1",
									"country",
									"locality",
									"postalcode",
									"region",
									"street",
								),
							),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				address: &domain.Address{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "user1",
					},
					Country:       "country",
					Locality:      "locality",
					PostalCode:    "postalcode",
					Region:        "region",
					StreetAddress: "street",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.Address{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "user1",
						ResourceOwner: "org1",
					},
					Country:       "country",
					Locality:      "locality",
					PostalCode:    "postalcode",
					Region:        "region",
					StreetAddress: "street",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeHumanAddress(tt.args.ctx, tt.args.address)
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

func newAddressChangedEvent(ctx context.Context, userID, resourceOwner, country, locality, postalCode, region, street string) *user.HumanAddressChangedEvent {
	event, _ := user.NewAddressChangedEvent(ctx,
		&user.NewAggregate(userID, resourceOwner).Aggregate,
		[]user.AddressChanges{
			user.ChangeCountry(country),
			user.ChangeLocality(locality),
			user.ChangePostalCode(postalCode),
			user.ChangeRegion(region),
			user.ChangeStreetAddress(street),
		},
	)
	return event
}
