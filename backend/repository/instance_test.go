package repository_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/cache/connector/gomap"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/database/mock"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

func Test_instance_Create(t *testing.T) {
	type args struct {
		ctx      context.Context
		tx       database.Transaction
		instance *repository.Instance
	}
	tests := []struct {
		name    string
		opts    []repository.Option[repository.InstanceOptions]
		args    args
		want    *repository.Instance
		wantErr bool
	}{
		{
			name: "simple",
			opts: []repository.Option[repository.InstanceOptions]{
				repository.WithTracer[repository.InstanceOptions](tracing.NewTracer("test")),
				repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
				repository.WithInstanceCache(
					repository.NewInstanceCache(gomap.NewCache[repository.InstanceIndex, string, *repository.Instance](context.Background(), repository.InstanceIndices, cache.Config{})),
				),
			},
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(t, mock.ExpectExec(repository.InstanceCreateStmt, "ID", "Name")),
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
		{
			name: "without cache",
			opts: []repository.Option[repository.InstanceOptions]{
				repository.WithTracer[repository.InstanceOptions](tracing.NewTracer("test")),
				repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
			},
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(t, mock.ExpectExec(repository.InstanceCreateStmt, "ID", "Name")),
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
		{
			name: "without cache, tracer",
			opts: []repository.Option[repository.InstanceOptions]{
				repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
			},
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(t, mock.ExpectExec(repository.InstanceCreateStmt, "ID", "Name")),
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
		{
			name: "without cache, tracer, logger",
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(t, mock.ExpectExec(repository.InstanceCreateStmt, "ID", "Name")),
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
		{
			name: "without cache, tracer, logger, eventStore",
			args: args{
				ctx: context.Background(),
				tx:  mock.NewTransaction(t, mock.ExpectExec(repository.InstanceCreateStmt, "ID", "Name")),
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
			fmt.Printf("------------------------ %s ------------------------\n", tt.name)
			i := repository.NewInstance(tt.opts...)
			got, err := i.Create(tt.args.ctx, tt.args.tx, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("instance.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instance.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_instance_ByID(t *testing.T) {
	type args struct {
		ctx context.Context
		tx  database.Transaction
		id  string
	}
	tests := []struct {
		name    string
		opts    []repository.Option[repository.InstanceOptions]
		args    args
		want    *repository.Instance
		wantErr bool
	}{
		{
			name: "simple, not cached",
			opts: []repository.Option[repository.InstanceOptions]{
				repository.WithTracer[repository.InstanceOptions](tracing.NewTracer("test")),
				repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
				repository.WithInstanceCache(
					repository.NewInstanceCache(gomap.NewCache[repository.InstanceIndex, string, *repository.Instance](context.Background(), repository.InstanceIndices, cache.Config{})),
				),
			},
			args: args{
				ctx: context.Background(),
				tx: mock.NewTransaction(t,
					mock.ExpectQueryRow(mock.NewRow(t, "id", "Name"), repository.InstanceByIDStmt, "id"),
				),
				id: "id",
			},
			want: &repository.Instance{
				ID:   "id",
				Name: "Name",
			},
			wantErr: false,
		},
		{
			name: "simple, cached",
			opts: []repository.Option[repository.InstanceOptions]{
				repository.WithTracer[repository.InstanceOptions](tracing.NewTracer("test")),
				repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
				repository.WithInstanceCache(
					func() *repository.InstanceCache {
						c := repository.NewInstanceCache(gomap.NewCache[repository.InstanceIndex, string, *repository.Instance](context.Background(), repository.InstanceIndices, cache.Config{}))
						c.Set(context.Background(), &repository.Instance{
							ID:   "id",
							Name: "Name",
						})
						return c
					}(),
				),
			},
			args: args{
				ctx: context.Background(),
				tx: mock.NewTransaction(t,
					mock.ExpectQueryRow(mock.NewRow(t, "id", "Name"), repository.InstanceByIDStmt, "id"),
				),
				id: "id",
			},
			want: &repository.Instance{
				ID:   "id",
				Name: "Name",
			},
			wantErr: false,
		},
		// {
		// 	name: "without cache, tracer",
		// 	opts: []repository.Option[repository.InstanceOptions]{
		// 		repository.WithLogger[repository.InstanceOptions](logging.New(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))),
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		tx:  mock.NewTransaction(),
		// 		id: &repository.Instance{
		// 			ID:   "ID",
		// 			Name: "Name",
		// 		},
		// 	},
		// 	want: &repository.Instance{
		// 		ID:   "ID",
		// 		Name: "Name",
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "without cache, tracer, logger",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		tx:  mock.NewTransaction(),
		// 		id: &repository.Instance{
		// 			ID:   "ID",
		// 			Name: "Name",
		// 		},
		// 	},
		// 	want: &repository.Instance{
		// 		ID:   "ID",
		// 		Name: "Name",
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("------------------------ %s ------------------------\n", tt.name)
			i := repository.NewInstance(tt.opts...)
			got, err := i.ByID(tt.args.ctx, tt.args.tx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("instance.ByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instance.ByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
