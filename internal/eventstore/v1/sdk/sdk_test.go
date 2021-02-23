package sdk

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestFilter(t *testing.T) {
	type args struct {
		filter   filterFunc
		appender appendFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr func(error) bool
	}{
		{
			name: "filter error",
			args: args{
				filter: func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error) {
					return nil, errors.ThrowInternal(nil, "test-46VX2", "test error")
				},
				appender: nil,
			},
			wantErr: errors.IsInternal,
		},
		{
			name: "no events found",
			args: args{
				filter: func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error) {
					return []*es_models.Event{}, nil
				},
				appender: nil,
			},
			wantErr: errors.IsNotFound,
		},
		{
			name: "append fails",
			args: args{
				filter: func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error) {
					return []*es_models.Event{&es_models.Event{}}, nil
				},
				appender: func(...*es_models.Event) error {
					return errors.ThrowInvalidArgument(nil, "SDK-DhBzl", "test error")
				},
			},
			wantErr: IsAppendEventError,
		},
		{
			name: "filter correct",
			args: args{
				filter: func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error) {
					return []*es_models.Event{&es_models.Event{}}, nil
				},
				appender: func(...*es_models.Event) error {
					return nil
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Filter(context.Background(), tt.args.filter, tt.args.appender, nil)
			if tt.wantErr == nil && err != nil {
				t.Errorf("no error expected %v", err)
			}
			if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("no error has wrong type %v", err)
			}
		})
	}
}

func TestPush(t *testing.T) {
	type args struct {
		push        pushFunc
		appender    appendFunc
		aggregaters []AggregateFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr func(error) bool
	}{
		{
			name: "no aggregates",
			args: args{
				push:        nil,
				appender:    nil,
				aggregaters: nil,
			},
			wantErr: errors.IsPreconditionFailed,
		},
		{
			name: "aggregater fails",
			args: args{
				push:     nil,
				appender: nil,
				aggregaters: []AggregateFunc{
					func(context.Context) (*es_models.Aggregate, error) {
						return nil, errors.ThrowInternal(nil, "SDK-Ec5x2", "test err")
					},
				},
			},
			wantErr: errors.IsInternal,
		},
		{
			name: "push fails",
			args: args{
				push: func(context.Context, ...*es_models.Aggregate) error {
					return errors.ThrowInternal(nil, "SDK-0g4gW", "test error")
				},
				appender: nil,
				aggregaters: []AggregateFunc{
					func(context.Context) (*es_models.Aggregate, error) {
						return &es_models.Aggregate{}, nil
					},
				},
			},
			wantErr: errors.IsInternal,
		},
		{
			name: "append aggregates fails",
			args: args{
				push: func(context.Context, ...*es_models.Aggregate) error {
					return nil
				},
				appender: func(...*es_models.Event) error {
					return errors.ThrowInvalidArgument(nil, "SDK-BDhcT", "test err")
				},
				aggregaters: []AggregateFunc{
					func(context.Context) (*es_models.Aggregate, error) {
						return &es_models.Aggregate{Events: []*es_models.Event{&es_models.Event{}}}, nil
					},
				},
			},
			wantErr: IsAppendEventError,
		},
		{
			name: "correct one aggregate",
			args: args{
				push: func(context.Context, ...*es_models.Aggregate) error {
					return nil
				},
				appender: func(...*es_models.Event) error {
					return nil
				},
				aggregaters: []AggregateFunc{
					func(context.Context) (*es_models.Aggregate, error) {
						return &es_models.Aggregate{Events: []*es_models.Event{&es_models.Event{}}}, nil
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "correct multiple aggregate",
			args: args{
				push: func(context.Context, ...*es_models.Aggregate) error {
					return nil
				},
				appender: func(...*es_models.Event) error {
					return nil
				},
				aggregaters: []AggregateFunc{
					func(context.Context) (*es_models.Aggregate, error) {
						return &es_models.Aggregate{Events: []*es_models.Event{&es_models.Event{}}}, nil
					},
					func(context.Context) (*es_models.Aggregate, error) {
						return &es_models.Aggregate{Events: []*es_models.Event{&es_models.Event{}}}, nil
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Push(context.Background(), tt.args.push, tt.args.appender, tt.args.aggregaters...)
			if tt.wantErr == nil && err != nil {
				t.Errorf("no error expected %v", err)
			}
			if tt.wantErr != nil && !tt.wantErr(err) {
				t.Errorf("no error has wrong type %v", err)
			}
		})
	}
}
