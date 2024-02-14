package execution

import (
	"reflect"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
)

func Test_createTargetToCommand(t *testing.T) {
	type args struct {
		req *execution.CreateTargetRequest
	}
	tests := []struct {
		name string
		args args
		want *command.Execution
	}{
		{
			name: "nil",
			args: args{nil},
			want: &command.Execution{
				Name:             "",
				ExecutionType:    domain.ExecutionTypeUndefined,
				URL:              "",
				Timeout:          0,
				Async:            false,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (async webhook)",
			args: args{&execution.CreateTargetRequest{
				Name:    "target 1",
				Type:    execution.TargetType_TARGET_TYPE_REST_WEBHOOK,
				Url:     "https://example.com/hooks/1",
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.CreateTargetRequest_IsAsync{
					IsAsync: true,
				},
			}},
			want: &command.Execution{
				Name:             "target 1",
				ExecutionType:    domain.ExecutionTypeWebhook,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            true,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&execution.CreateTargetRequest{
				Name:    "target 1",
				Type:    execution.TargetType_TARGET_TYPE_REST_REQUEST_RESPONSE,
				Url:     "https://example.com/hooks/1",
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.CreateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
			}},
			want: &command.Execution{
				Name:             "target 1",
				ExecutionType:    domain.ExecutionTypeRequestResponse,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            false,
				InterruptOnError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTargetToCommand(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTargetToCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateTargetToCommand(t *testing.T) {
	type args struct {
		req *execution.UpdateTargetRequest
	}
	tests := []struct {
		name string
		args args
		want *command.Execution
	}{
		{
			name: "nil",
			args: args{nil},
			want: &command.Execution{
				Name:             "",
				ExecutionType:    domain.ExecutionTypeUndefined,
				URL:              "",
				Timeout:          0,
				Async:            false,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields nil",
			args: args{&execution.UpdateTargetRequest{
				Name:          nil,
				Type:          nil,
				Url:           nil,
				Timeout:       nil,
				ExecutionType: nil,
			}},
			want: &command.Execution{
				Name:             "",
				ExecutionType:    domain.ExecutionTypeUndefined,
				URL:              "",
				Timeout:          0,
				Async:            false,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (async webhook)",
			args: args{&execution.UpdateTargetRequest{
				Name:    gu.Ptr("target 1"),
				Type:    gu.Ptr(execution.TargetType_TARGET_TYPE_REST_WEBHOOK),
				Url:     gu.Ptr("https://example.com/hooks/1"),
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.UpdateTargetRequest_IsAsync{
					IsAsync: true,
				},
			}},
			want: &command.Execution{
				Name:             "target 1",
				ExecutionType:    domain.ExecutionTypeWebhook,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            true,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&execution.UpdateTargetRequest{
				Name:    gu.Ptr("target 1"),
				Type:    gu.Ptr(execution.TargetType_TARGET_TYPE_REST_REQUEST_RESPONSE),
				Url:     gu.Ptr("https://example.com/hooks/1"),
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &execution.UpdateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
			}},
			want: &command.Execution{
				Name:             "target 1",
				ExecutionType:    domain.ExecutionTypeRequestResponse,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            false,
				InterruptOnError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateTargetToCommand(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTargetToCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
