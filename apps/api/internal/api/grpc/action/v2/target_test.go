package action

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
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
				Endpoint:         "",
				Timeout:          0,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (webhook)",
			args: args{&action.CreateTargetRequest{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.AddTarget{
				Name:             "target 1",
				TargetType:       domain.TargetTypeWebhook,
				Endpoint:         "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (async)",
			args: args{&action.CreateTargetRequest{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.CreateTargetRequest_RestAsync{
					RestAsync: &action.RESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.AddTarget{
				Name:             "target 1",
				TargetType:       domain.TargetTypeAsync,
				Endpoint:         "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
				InterruptOnError: false,
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&action.CreateTargetRequest{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.AddTarget{
				Name:             "target 1",
				TargetType:       domain.TargetTypeCall,
				Endpoint:         "https://example.com/hooks/1",
				Timeout:          10 * time.Second,
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
				Name:       nil,
				TargetType: nil,
				Timeout:    nil,
			}},
			want: &command.ChangeTarget{
				Name:             nil,
				TargetType:       nil,
				Endpoint:         nil,
				Timeout:          nil,
				InterruptOnError: nil,
			},
		},
		{
			name: "all fields empty",
			args: args{&action.UpdateTargetRequest{
				Name:       gu.Ptr(""),
				TargetType: nil,
				Timeout:    durationpb.New(0),
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr(""),
				TargetType:       nil,
				Endpoint:         nil,
				Timeout:          gu.Ptr(0 * time.Second),
				InterruptOnError: nil,
			},
		},
		{
			name: "all fields (webhook)",
			args: args{&action.UpdateTargetRequest{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.UpdateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeWebhook),
				Endpoint:         gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
				InterruptOnError: gu.Ptr(false),
			},
		},
		{
			name: "all fields (webhook interrupt)",
			args: args{&action.UpdateTargetRequest{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.UpdateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeWebhook),
				Endpoint:         gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
				InterruptOnError: gu.Ptr(true),
			},
		},
		{
			name: "all fields (async)",
			args: args{&action.UpdateTargetRequest{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.UpdateTargetRequest_RestAsync{
					RestAsync: &action.RESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeAsync),
				Endpoint:         gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
				InterruptOnError: gu.Ptr(false),
			},
		},
		{
			name: "all fields (interrupting response)",
			args: args{&action.UpdateTargetRequest{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.UpdateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			}},
			want: &command.ChangeTarget{
				Name:             gu.Ptr("target 1"),
				TargetType:       gu.Ptr(domain.TargetTypeCall),
				Endpoint:         gu.Ptr("https://example.com/hooks/1"),
				Timeout:          gu.Ptr(10 * time.Second),
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
