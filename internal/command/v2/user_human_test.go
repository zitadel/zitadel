package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestAddHumanCommand(t *testing.T) {
	type args struct {
		a           *user.Aggregate
		human       *AddHuman
		passwordAlg crypto.HashAlgorithm
		filter      preparation.FilterToQueryReducer
	}
	agg := user.NewAggregate("id", "ro")
	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid email",
			args: args{
				a: agg,
				human: &AddHuman{
					Email: "invalid",
				},
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "USER-Ec7dM", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid first name",
			args: args{
				a: agg,
				human: &AddHuman{
					Email: "support@zitadel.ch",
				},
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "USER-UCej2", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid last name",
			args: args{
				a: agg,
				human: &AddHuman{
					Email:     "support@zitadel.ch",
					FirstName: "hurst",
				},
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "USER-DiAq8", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid password",
			args: args{
				a: agg,
				human: &AddHuman{
					Email:     "support@zitadel.ch",
					FirstName: "gigi",
					LastName:  "giraffe",
					Password:  "short",
				},
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							org.NewDomainPolicyAddedEvent(
								context.Background(),
								&org.NewAggregate("id", "ro").Aggregate,
								true,
							),
						}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									context.Background(),
									&org.NewAggregate("id", "ro").Aggregate,
									8,
									true,
									true,
									true,
									true,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowInvalidArgument(nil, "COMMA-HuJf6", "Errors.User.PasswordComplexityPolicy.MinLength"),
			},
		},
		{
			name: "correct",
			args: args{
				a: agg,
				human: &AddHuman{
					Email:     "support@zitadel.ch",
					FirstName: "gigi",
					LastName:  "giraffe",
					Password:  "",
				},
				passwordAlg: crypto.CreateMockHashAlg(gomock.NewController(t)),
				filter: NewMultiFilter().Append(
					func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							org.NewDomainPolicyAddedEvent(
								context.Background(),
								&org.NewAggregate("id", "ro").Aggregate,
								true,
							),
						}, nil
					}).
					Append(
						func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
							return []eventstore.Event{
								org.NewPasswordComplexityPolicyAddedEvent(
									context.Background(),
									&org.NewAggregate("id", "ro").Aggregate,
									2,
									false,
									false,
									false,
									false,
								),
							}, nil
						}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					user.NewHumanAddedEvent(
						context.Background(),
						&agg.Aggregate,
						"",
						"gigi",
						"giraffe",
						"",
						"",
						language.Und,
						0,
						"support@zitadel.ch",
						true,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, AddHumanCommand(tt.args.a, tt.args.human, tt.args.passwordAlg), tt.args.filter, tt.want)
		})
	}
}
