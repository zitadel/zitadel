package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
)

func Test_sessionsToPb(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)

	sessions := []*query.Session{
		{ // no factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			Metadata:      map[string][]byte{"hello": []byte("world")},
		},
		{ // user factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			UserFactor: query.SessionUserFactor{
				UserID:        "345",
				UserCheckedAt: past,
				LoginName:     "donald",
				DisplayName:   "donald duck",
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // no factor
			ID:            "999",
			CreationDate:  now,
			ChangeDate:    now,
			Sequence:      123,
			State:         domain.SessionStateActive,
			ResourceOwner: "me",
			Creator:       "he",
			PasswordFactor: query.SessionPasswordFactor{
				PasswordCheckedAt: past,
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
	}

	want := []*session.Session{
		{ // no factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors:      nil,
			Metadata:     map[string][]byte{"hello": []byte("world")},
		},
		{ // user factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors: &session.Factors{
				User: &session.UserFactor{
					VerifiedAt:  timestamppb.New(past),
					Id:          "345",
					LoginName:   "donald",
					DisplayName: "donald duck",
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
		{ // password factor
			Id:           "999",
			CreationDate: timestamppb.New(now),
			ChangeDate:   timestamppb.New(now),
			Sequence:     123,
			Factors: &session.Factors{
				Password: &session.PasswordFactor{
					VerifiedAt: timestamppb.New(past),
				},
			},
			Metadata: map[string][]byte{"hello": []byte("world")},
		},
	}

	out := sessionsToPb(sessions)
	require.Len(t, out, len(want))

	for i, got := range out {
		if !proto.Equal(got, want[i]) {
			t.Errorf("session %d got:\n%v\nwant:\n%v", i, got, want)
		}
	}
}

