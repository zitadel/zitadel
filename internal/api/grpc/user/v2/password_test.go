package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/domain"
)

func Test_notificationTypeToDomain(t *testing.T) {
	tests := []struct {
		name             string
		notificationType user.NotificationType
		want             domain.NotificationType
	}{
		{
			"unspecified",
			user.NotificationType_NOTIFICATION_TYPE_Unspecified,
			domain.NotificationTypeEmail,
		},
		{
			"email",
			user.NotificationType_NOTIFICATION_TYPE_Email,
			domain.NotificationTypeEmail,
		},
		{
			"sms",
			user.NotificationType_NOTIFICATION_TYPE_SMS,
			domain.NotificationTypeSms,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, notificationTypeToDomain(tt.notificationType), "notificationTypeToDomain(%v)", tt.notificationType)
		})
	}
}
