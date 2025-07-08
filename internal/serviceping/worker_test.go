package serviceping

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/serviceping/mock"
	"github.com/zitadel/zitadel/internal/zerrors"
	analytics "github.com/zitadel/zitadel/pkg/grpc/analytics/v2beta"
)

var (
	testNow   = time.Now()
	errInsert = fmt.Errorf("insert error")
)

func TestWorker_reportBaseInformation(t *testing.T) {
	type fields struct {
		reportClient func(*testing.T) analytics.TelemetryServiceClient
		db           func(*testing.T) Queries
		systemID     string
		version      string
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		reportID string
		err      error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "database error, error",
			fields: fields{
				db: func(t *testing.T) Queries {
					ctrl := gomock.NewController(t)
					queries := mock.NewMockQueries(ctrl)
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						nil, zerrors.ThrowInternal(nil, "id", "db error"),
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					return mock.NewMockTelemetryServiceClient(gomock.NewController(t))
				},
			},
			want: want{
				reportID: "",
				err:      zerrors.ThrowInternal(nil, "id", "db error"),
			},
		},
		{
			name: "telemetry client error, error",
			fields: fields{
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain", "domain2"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						nil, status.Error(codes.Internal, "error"),
					)
					return client
				},
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
										{
											Domain: "domain2",
										},
									},
								},
							},
						},
						nil,
					)
					return queries
				},
				systemID: "system-id",
				version:  "version",
			},
			want: want{
				reportID: "",
				err:      status.Error(codes.Internal, "error"),
			},
		},
		{
			name: "report ok, reportID returned",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
										{
											Domain: "domain2",
										},
									},
								},
							},
						},
						nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain", "domain2"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						&analytics.ReportBaseInformationResponse{ReportId: "report-id"}, nil,
					)
					return client
				},
				systemID: "system-id",
				version:  "version",
			},
			want: want{
				reportID: "report-id",
				err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				reportClient: tt.fields.reportClient(t),
				db:           tt.fields.db(t),
				systemID:     tt.fields.systemID,
				version:      tt.fields.version,
			}
			got, err := w.reportBaseInformation(tt.args.ctx)
			assert.Equal(t, tt.want.reportID, got)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestWorker_reportResourceCounts(t *testing.T) {
	type fields struct {
		reportClient func(*testing.T) analytics.TelemetryServiceClient
		db           func(*testing.T) Queries
		config       *Config
		systemID     string
	}
	type args struct {
		ctx      context.Context
		reportID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "database error, error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 1).Return(
						nil, zerrors.ThrowInternal(nil, "id", "db error"),
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					return mock.NewMockTelemetryServiceClient(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 1,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx:      context.Background(),
				reportID: "",
			},
			wantErr: zerrors.ThrowInternal(nil, "id", "db error"),
		},
		{
			name: "no resource counts, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 1).Return(
						[]query.ResourceCount{}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					return mock.NewMockTelemetryServiceClient(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 1,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx:      context.Background(),
				reportID: "",
			},
			wantErr: nil,
		},
		{
			name: "telemetry client error, error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: nil,
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
						},
					}).Return(
						nil, status.Error(codes.Internal, "error"),
					)
					return client
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx:      context.Background(),
				reportID: "",
			},
			wantErr: status.Error(codes.Internal, "error"),
		},
		{
			name: "report ok, no additional counts, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: nil,
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
						},
					}).Return(
						&analytics.ReportResourceCountsResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx:      context.Background(),
				reportID: "",
			},
			wantErr: nil,
		},
		{
			name: "report ok, additional counts, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
							{
								ID:         2,
								InstanceID: "instance-id2",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id2",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     5,
							},
						}, nil,
					)
					queries.EXPECT().ListResourceCounts(gomock.Any(), 2, 2).Return(
						[]query.ResourceCount{
							{
								ID:         3,
								InstanceID: "instance-id3",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id3",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     20,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: nil,
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
							{
								InstanceId:   "instance-id2",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id2",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       5,
							},
						},
					}).Return(
						&analytics.ReportResourceCountsResponse{
							ReportId: "report-id",
						}, nil,
					)
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: gu.Ptr("report-id"),
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id3",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id3",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       20,
							},
						},
					}).Return(
						&analytics.ReportResourceCountsResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx:      context.Background(),
				reportID: "",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				reportClient: tt.fields.reportClient(t),
				db:           tt.fields.db(t),
				config:       tt.fields.config,
				systemID:     tt.fields.systemID,
			}
			err := w.reportResourceCounts(tt.args.ctx, tt.args.reportID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWorker_Work(t *testing.T) {
	type fields struct {
		WorkerDefaults river.WorkerDefaults[*ServicePingReport]
		reportClient   func(*testing.T) analytics.TelemetryServiceClient
		db             func(*testing.T) Queries
		queue          func(*testing.T) Queue
		config         *Config
		systemID       string
		version        string
	}
	type args struct {
		ctx context.Context
		job *river.Job[*ServicePingReport]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "unknown report type, cancel job",
			fields: fields{
				db: func(t *testing.T) Queries {
					return mock.NewMockQueries(gomock.NewController(t))
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					return mock.NewMockTelemetryServiceClient(gomock.NewController(t))
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: 100000,
					},
				},
			},
			wantErr: river.JobCancel(ErrInvalidReportType),
		},
		{
			name: "report base information, database error, retry job",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						nil, zerrors.ThrowInternal(nil, "id", "db error"),
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					return mock.NewMockTelemetryServiceClient(gomock.NewController(t))
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeBaseInformation,
					},
				},
			},
			wantErr: zerrors.ThrowInternal(nil, "id", "db error"),
		},
		{
			name: "report base information, config error, cancel job",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
									},
								},
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						nil, &TelemetryError{StatusCode: http.StatusNotFound, Body: []byte("endpoint not found")},
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
				systemID: "system-id",
				version:  "version",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeBaseInformation,
					},
				},
			},
			wantErr: river.JobCancel(&TelemetryError{StatusCode: http.StatusNotFound, Body: []byte("endpoint not found")}),
		},
		{
			name: "report base information, no reports enabled, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
									},
								},
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						&analytics.ReportBaseInformationResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							Enabled: false,
						},
					},
				},
				systemID: "system-id",
				version:  "version",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeBaseInformation,
					},
				},
			},
		},
		{
			name: "report base information, job creation error, cancel job",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
									},
								},
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						&analytics.ReportBaseInformationResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					q := mock.NewMockQueue(gomock.NewController(t))
					q.EXPECT().Insert(gomock.Any(),
						&ServicePingReport{
							ReportID:   "report-id",
							ReportType: ReportTypeResourceCounts,
						},
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithQueueName(QueueName))),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithMaxAttempts(5)))). // TODO: better solution
						Return(errInsert)
					return q
				},
				config: &Config{
					MaxAttempts: 5,
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							Enabled: true,
						},
					},
				},
				systemID: "system-id",
				version:  "version",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeBaseInformation,
					},
				},
			},
			wantErr: errInsert,
		},
		{
			name: "report base information, success, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().SearchInstances(gomock.Any(), &query.InstanceSearchQueries{}).Return(
						&query.Instances{
							Instances: []*query.Instance{
								{
									ID:           "id",
									CreationDate: testNow,
									Domains: []*query.InstanceDomain{
										{
											Domain: "domain",
										},
									},
								},
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportBaseInformation(gomock.Any(), &analytics.ReportBaseInformationRequest{
						SystemId: "system-id",
						Version:  "version",
						Instances: []*analytics.InstanceInformation{
							{
								Id:        "id",
								Domains:   []string{"domain"},
								CreatedAt: timestamppb.New(testNow),
							},
						},
					}).Return(
						&analytics.ReportBaseInformationResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					q := mock.NewMockQueue(gomock.NewController(t))
					q.EXPECT().Insert(gomock.Any(),
						&ServicePingReport{
							ReportID:   "report-id",
							ReportType: ReportTypeResourceCounts,
						},
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithQueueName(QueueName))),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithMaxAttempts(5)))).
						Return(nil)
					return q
				},
				config: &Config{
					MaxAttempts: 5,
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							Enabled: true,
						},
					},
				},
				systemID: "system-id",
				version:  "version",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeBaseInformation,
					},
				},
			},
		},
		{
			name: "report resource counts, service unavailable, retry job",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: gu.Ptr("report-id"),
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
						},
					}).Return(
						nil, status.Error(codes.Unavailable, "service unavailable"),
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportType: ReportTypeResourceCounts,
						ReportID:   "report-id",
					},
				},
			},
			wantErr: status.Error(codes.Unavailable, "service unavailable"),
		},
		{
			name: "report resource counts, precondition error, cancel job",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: gu.Ptr("report-id"),
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
						},
					}).Return(
						nil, &TelemetryError{StatusCode: http.StatusPreconditionFailed, Body: []byte("report too old")},
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportID:   "report-id",
						ReportType: ReportTypeResourceCounts,
					},
				},
			},
			wantErr: river.JobCancel(&TelemetryError{StatusCode: http.StatusPreconditionFailed, Body: []byte("report too old")}),
		},
		{
			name: "report resource counts, success, no error",
			fields: fields{
				db: func(t *testing.T) Queries {
					queries := mock.NewMockQueries(gomock.NewController(t))
					queries.EXPECT().ListResourceCounts(gomock.Any(), 0, 2).Return(
						[]query.ResourceCount{
							{
								ID:         1,
								InstanceID: "instance-id",
								TableName:  "table_name",
								ParentType: domain.CountParentTypeInstance,
								ParentID:   "instance-id",
								Resource:   "resource",
								UpdatedAt:  testNow,
								Amount:     10,
							},
						}, nil,
					)
					return queries
				},
				reportClient: func(t *testing.T) analytics.TelemetryServiceClient {
					client := mock.NewMockTelemetryServiceClient(gomock.NewController(t))
					client.EXPECT().ReportResourceCounts(gomock.Any(), &analytics.ReportResourceCountsRequest{
						SystemId: "system-id",
						ReportId: gu.Ptr("report-id"),
						ResourceCounts: []*analytics.ResourceCount{
							{
								InstanceId:   "instance-id",
								TableName:    "table_name",
								ParentType:   analytics.CountParentType_COUNT_PARENT_TYPE_INSTANCE,
								ParentId:     "instance-id",
								ResourceName: "resource",
								UpdatedAt:    timestamppb.New(testNow),
								Amount:       10,
							},
						},
					}).Return(
						&analytics.ReportResourceCountsResponse{
							ReportId: "report-id",
						}, nil,
					)
					return client
				},
				queue: func(t *testing.T) Queue {
					return mock.NewMockQueue(gomock.NewController(t))
				},
				config: &Config{
					Telemetry: TelemetryConfig{
						ResourceCount: ResourceCount{
							BulkSize: 2,
						},
					},
				},
				systemID: "system-id",
			},
			args: args{
				ctx: context.Background(),
				job: &river.Job[*ServicePingReport]{
					Args: &ServicePingReport{
						ReportID:   "report-id",
						ReportType: ReportTypeResourceCounts,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				WorkerDefaults: river.WorkerDefaults[*ServicePingReport]{},
				reportClient:   tt.fields.reportClient(t),
				db:             tt.fields.db(t),
				queue:          tt.fields.queue(t),
				config:         tt.fields.config,
				systemID:       tt.fields.systemID,
				version:        tt.fields.version,
			}
			err := w.Work(tt.args.ctx, tt.args.job)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func Test_parseAndValidateSchedule(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name          string
		args          args
		wantNextStart time.Time
		wantNextEnd   time.Time
		wantErr       error
	}{
		{
			name: "@daily, returns randomized daily schedule",
			args: args{
				interval: "@daily",
			},
			wantNextStart: time.Now(),
			wantNextEnd:   time.Now().Add(24 * time.Hour),
		},
		{
			name: "invalid cron expression, returns error",
			args: args{
				interval: "invalid cron",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "SERV-NJqiof", "invalid interval"),
		},
		{
			name: "valid cron expression, returns schedule",
			args: args{
				interval: "0 0 * * *",
			},
			wantNextStart: nextMidnight(),
			wantNextEnd:   nextMidnight(),
		},
		{
			name: "valid cron expression (extended syntax), returns schedule",
			args: args{
				interval: "@midnight",
			},
			wantNextStart: nextMidnight(),
			wantNextEnd:   nextMidnight(),
		},
		{
			name: "less than minInterval, returns error",
			args: args{
				interval: "0/15 * * * *",
			},
			wantErr: zerrors.ThrowInvalidArgumentf(nil, "SERV-FJ12", "interval must be at least %s", minInterval),
		},
		{
			name: "less than minInterval (extended syntax), returns error",
			args: args{
				interval: "@every 15m",
			},
			wantErr: zerrors.ThrowInvalidArgumentf(nil, "SERV-FJ12", "interval must be at least %s", minInterval),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAndValidateSchedule(tt.args.interval)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				now := time.Now()
				assert.WithinRange(t, got.Next(now), tt.wantNextStart, tt.wantNextEnd)
			}
		})
	}
}

func nextMidnight() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, time.Local)
}
