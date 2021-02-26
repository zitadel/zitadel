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
