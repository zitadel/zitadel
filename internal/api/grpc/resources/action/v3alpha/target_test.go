package action

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/internal/domain"
	action "github.com/zitadel/zitadel/v2/pkg/grpc/resources/action/v3alpha"
)

func Test_createTargetToCommand(t *testing.T) {
	type args struct {
		req *action.Target
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
			args: args{&action.Target{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.Target_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{},
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
			args: args{&action.Target{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.Target_RestAsync{
					RestAsync: &action.SetRESTAsync{},
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
			args: args{&action.Target{
				Name:     "target 1",
				Endpoint: "https://example.com/hooks/1",
				TargetType: &action.Target_RestCall{
					RestCall: &action.SetRESTCall{
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
			got := createTargetToCommand(&action.CreateTargetRequest{Target: tt.args.req})
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_updateTargetToCommand(t *testing.T) {
	type args struct {
		req *action.PatchTarget
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
			args: args{&action.PatchTarget{
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
			args: args{&action.PatchTarget{
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
			args: args{&action.PatchTarget{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.PatchTarget_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
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
			args: args{&action.PatchTarget{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.PatchTarget_RestWebhook{
					RestWebhook: &action.SetRESTWebhook{
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
			args: args{&action.PatchTarget{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.PatchTarget_RestAsync{
					RestAsync: &action.SetRESTAsync{},
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
			args: args{&action.PatchTarget{
				Name:     gu.Ptr("target 1"),
				Endpoint: gu.Ptr("https://example.com/hooks/1"),
				TargetType: &action.PatchTarget_RestCall{
					RestCall: &action.SetRESTCall{
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
			got := patchTargetToCommand(&action.PatchTargetRequest{Target: tt.args.req})
			assert.Equal(t, tt.want, got)
		})
	}
}
