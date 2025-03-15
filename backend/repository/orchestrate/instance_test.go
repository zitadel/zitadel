package orchestrate

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/cache"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/mock"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

func Test_instance_SetUp(t *testing.T) {
	type fields struct {
		options options
		cache   *cache.Instance
	}
	type args struct {
		ctx      context.Context
		tx       database.Transaction
		instance *repository.Instance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *repository.Instance
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{
				options: options{
					tracer: tracing.NewTracer("test"),
					logger: logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))),
				},
				cache: cache.NewInstance(),
			},
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(),
				instance: &repository.Instance{
					ID:   "ID",
					Name: "Name",
				},
			},
			want: &repository.Instance{
				ID:   "ID",
				Name: "Name",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &instance{
				options: tt.fields.options,
				cache:   tt.fields.cache,
			}
			got, err := i.SetUp(tt.args.ctx, tt.args.tx, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("instance.SetUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instance.SetUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
