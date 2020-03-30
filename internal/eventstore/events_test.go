package eventstore

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

func Test_app_CreateEvents(t *testing.T) {
	type fields struct {
		repo repository.Repository
	}
	type args struct {
		ctx        context.Context
		aggregates []*models.Aggregate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &app{
				repo: tt.fields.repo,
			}
			if err := es.CreateEvents(tt.args.ctx, tt.args.aggregates...); (err != nil) != tt.wantErr {
				t.Errorf("app.CreateEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
