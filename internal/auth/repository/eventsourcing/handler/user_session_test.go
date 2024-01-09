package handler

import (
	"testing"

	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

// tests the proper working of the cache function
func TestUserSession_fillUserInfo(t *testing.T) {
	type args struct {
		sessions []*view_model.UserSessionView
	}
	tests := []struct {
		name      string
		args      args
		cacheHits map[string]int
	}{
		{
			name: "one session",
			args: args{
				sessions: []*view_model.UserSessionView{
					{
						UserID:     "user",
						InstanceID: "instance",
					},
				},
			},
			cacheHits: map[string]int{
				"user-instance": 1,
			},
		},
		{
			name: "same user",
			args: args{
				sessions: []*view_model.UserSessionView{
					{
						UserID:     "user",
						InstanceID: "instance",
					},
					{
						UserID:     "user",
						InstanceID: "instance",
					},
				},
			},
			cacheHits: map[string]int{
				"user-instance": 2,
			},
		},
		{
			name: "different users",
			args: args{
				sessions: []*view_model.UserSessionView{
					{
						UserID:     "user",
						InstanceID: "instance",
					},
					{
						UserID:     "user2",
						InstanceID: "instance",
					},
				},
			},
			cacheHits: map[string]int{
				"user-instance":  1,
				"user2-instance": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := map[string]int{}
			getUserByID := func(userID, instanceID string) (*view_model.UserView, error) {
				cache[userID+"-"+instanceID]++
				return &view_model.UserView{HumanView: &view_model.HumanView{}}, nil
			}
			for _, session := range tt.args.sessions {
				if err := new(UserSession).fillUserInfo(session, getUserByID); err != nil {
					t.Errorf("UserSession.fillUserInfo() unexpected error = %v", err)
				}
			}
			if len(cache) != len(tt.cacheHits) {
				t.Errorf("unexpected length of cache hits: want %d, got %d", len(tt.cacheHits), len(cache))
				return
			}
			for key, count := range tt.cacheHits {
				if cache[key] != count {
					t.Errorf("unexpected cache hits on %s: want %d, got %d", key, count, cache[key])
				}
			}
		})
	}
}
