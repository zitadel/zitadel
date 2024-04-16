package action

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
)

func Test_createTargetToCommand(t *testing.T) {
	type args struct {
		req *action.CreateTargetRequest
	}
	tests := []struct {
		name string
		args args
		want *command.AddTarget
	}{
		{
			name: "nil",
			args: args{nil},
			want: &command.AddTarget{
				Name:             "",
				URL:              "",
				Timeout:          0,
				Async:            false,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (async webhook)",
			args: args{&action.CreateTargetRequest{
				Name: "target 1",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com/hooks/1",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.CreateTargetRequest_IsAsync{
					IsAsync: true,
				},
			}},
			want: &command.AddTarget{
				Name:             "target 1",
				TargetType:       domain.TargetTypeWebhook,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            true,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&action.CreateTargetRequest{
				Name: "target 1",
				TargetType: &action.CreateTargetRequest_RestRequestResponse{
					RestRequestResponse: &action.SetRESTRequestResponse{
						Url: "https://example.com/hooks/1",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.CreateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
			}},
			want: &command.AddTarget{
				Name:             "target 1",
				TargetType:       domain.TargetTypeRequestResponse,
				URL:              "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				Async:            false,
				InterruptOnError: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createTargetToCommand(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_updateTargetToCommand(t *testing.T) {
	type args struct {
		req *action.UpdateTargetRequest
	}
	tests := []struct {
		name string
		args args
		want *command.ChangeTarget
	}{
		{
			name: "nil",
			args: args{nil},
			want: nil,
		},
		{
			name: "all fields nil",
			args: args{&action.UpdateTargetRequest{
				Name:          nil,
				TargetType:    nil,
				Timeout:       nil,
				ExecutionType: nil,
			}},
			want: &command.ChangeTarget{
				Name:             nil,
				TargetType:       nil,
				URL:              nil,
				Timeout:          nil,
				Async:            nil,
				InterruptOnError: nil,
			},
		},
		{
			name: "all fields empty",
			args: args{&action.UpdateTargetRequest{
				Name:          gu.Ptr(""),
				TargetType:    nil,
				Timeout:       durationpb.New(0),
				ExecutionType: nil,
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr(""),
				TargetType:       nil,
				URL:              nil,
				Timeout:          gu.Ptr(0 * time.Second),
				Async:            nil,
				InterruptOnError: nil,
			},
		},
		{
			name: "all fields (async webhook)",
			args: args{&action.UpdateTargetRequest{
				Name: gu.Ptr("target 1"),
				TargetType: &action.UpdateTargetRequest_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
						Url: "https://example.com/hooks/1",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.UpdateTargetRequest_IsAsync{
					IsAsync: true,
				},
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeWebhook),
				URL:              gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
				Async:            gu.Ptr(true),
				InterruptOnError: gu.Ptr(false),
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&action.UpdateTargetRequest{
				Name: gu.Ptr("target 1"),
				TargetType: &action.UpdateTargetRequest_RestRequestResponse{
					RestRequestResponse: &action.SetRESTRequestResponse{
						Url: "https://example.com/hooks/1",
					},
				},
				Timeout: durationpb.New(10 * time.Second),
				ExecutionType: &action.UpdateTargetRequest_InterruptOnError{
					InterruptOnError: true,
				},
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeRequestResponse),
				URL:              gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
				Async:            gu.Ptr(false),
				InterruptOnError: gu.Ptr(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateTargetToCommand(tt.args.req)
			assert.Equal(t, tt.want, got)
		})
	}
}
