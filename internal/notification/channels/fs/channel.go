package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/k3a/html2text"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func InitFSChannel[T channels.Message](config Config) (channels.NotificationChannel[T], error) {
	if err := os.MkdirAll(config.Path, os.ModePerm); err != nil {
		return nil, err
	}

	logging.Debug("successfully initialized filesystem email and sms channel")

	return channels.HandleMessageFunc[T](func(message T) error {
		fileName := fmt.Sprintf("%d_", time.Now().Unix())
		content, err := message.GetContent()
		if err != nil {
			return err
		}
		switch msg := any(message).(type) {
		case *messages.Email:
			recipients := make([]string, len(msg.Recipients))
			copy(recipients, msg.Recipients)
			sort.Strings(recipients)
			fileName = fileName + "mail_to_" + strings.Join(recipients, "_") + ".html"
			if config.Compact {
				content = html2text.HTML2Text(content)
			}
		case *messages.SMS:
			fileName = fileName + "sms_to_" + msg.RecipientPhoneNumber + ".txt"
		case *messages.JSON:
			fileName = "message.json"
		default:
			return errors.ThrowUnimplementedf(nil, "NOTIF-6f9a1", "filesystem provider doesn't support message type %T", message)
		}

		return os.WriteFile(filepath.Join(config.Path, fileName), []byte(content), 0666)
	}), nil
}
