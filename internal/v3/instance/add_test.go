package instance_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/zitadel/zitadel/internal/v3/instance"
	repo_log "github.com/zitadel/zitadel/internal/v3/repository/log"
	repo_mem "github.com/zitadel/zitadel/internal/v3/repository/memory"
	"github.com/zitadel/zitadel/internal/v3/storage"
	"github.com/zitadel/zitadel/internal/v3/storage/memory"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestBusinessLogic_AddInstance(t *testing.T) {
	type fields struct {
		client storage.Client
		stores []instance.InstanceStorage
	}
	type args struct {
		ctx     context.Context
		request *instance.AddInstanceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				client: &memory.Client{},
				stores: []instance.InstanceStorage{
					repo_mem.NewInstanceMemory(),
					repo_log.NewInstanceLogger(slog.Default()),
				},
			},
			args: args{
				ctx: context.Background(),
				request: &instance.AddInstanceRequest{
					AddInstanceRequest: system_pb.AddInstanceRequest{
						InstanceName:    "test",
						CustomDomain:    "test",
						DefaultLanguage: "en",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bl := instance.NewBusinessLogic(tt.fields.client, tt.fields.stores...)
			_, err := bl.AddInstance(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BusinessLogic.AddInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
