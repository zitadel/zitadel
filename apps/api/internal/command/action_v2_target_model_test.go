package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/target"
)

func targetAddEvent(aggID, resourceOwner string) *target.AddedEvent {
	return target.NewAddedEvent(context.Background(),
		target.NewAggregate(aggID, resourceOwner),
		"name",
		domain.TargetTypeWebhook,
		"https://example.com",
		time.Second,
		false,
		&crypto.CryptoValue{
			CryptoType: crypto.TypeEncryption,
			Algorithm:  "enc",
			KeyID:      "id",
			Crypted:    []byte("12345678"),
		},
	)
}

func targetRemoveEvent(aggID, resourceOwner string) *target.RemovedEvent {
	return target.NewRemovedEvent(context.Background(),
		target.NewAggregate(aggID, resourceOwner),
		"name",
	)
}

func TestCommandSide_targetsExistsWriteModel(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		ids           []string
		resourceOwner string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		res    bool
	}{
		{
			name: "target, single",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target"},
			},
			res: true,
		},
		{
			name: "target, single reset",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target", "org1"),
						),
						eventFromEventPusher(

							targetAddEvent("target", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target"},
			},
			res: true,
		},
		{
			name: "target, single before removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetRemoveEvent("target", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target"},
			},
			res: true,
		},
		{
			name: "target, single removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target"},
			},
			res: false,
		},
		{
			name: "target, multiple",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: true,
		},
		{
			name: "target, multiple, first removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: false,
		},
		{
			name: "target, multiple, second removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: false,
		},
		{
			name: "target, multiple, third removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: false,
		},
		{
			name: "target, multiple, before removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetRemoveEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: true,
		},
		{
			name: "target, multiple, all removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target3", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: false,
		},

		{
			name: "target, multiple, two removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							targetAddEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target1", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target2", "org1"),
						),
						eventFromEventPusher(
							targetAddEvent("target3", "org1"),
						),
						eventFromEventPusher(
							targetRemoveEvent("target2", "org1"),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"target1", "target2", "target3"},
			},
			res: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			assert.Equal(t, tt.res, c.existsTargetsByIDs(tt.args.ctx, tt.args.ids, tt.args.resourceOwner))

		})
	}
}
