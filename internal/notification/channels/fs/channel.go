package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/k3a/html2text"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/messages"
)

var _ channels.NotificationChannel = (*FS)(nil)

type FS struct {
	Path    string
	Compact bool
}

func InitFSProvider(config FSConfig) (*FS, error) {

	if err := os.MkdirAll(config.Path, os.ModePerm); err != nil {
		return nil, err
	}

	return &FS{
		Path: config.Path,
	}, nil
}

func (f *FS) HandleMessage(message channels.Message) error {

	var (
		fileName string
		content  = message.GetContent()
	)
	switch msg := message.(type) {
	case *messages.Email:
		recipients := make([]string, len(msg.Recipients))
		copy(recipients, msg.Recipients)
		sort.Strings(recipients)
		fileName = "mail_to_" + strings.Join(recipients, "_") + ".html"
		if f.Compact {
			content = html2text.HTML2Text(content)
		}
	case *messages.SMS:
		fileName = "sms_to_" + msg.RecipientPhoneNumber + ".txt"
	default:
		logging.Log("NOTIF-6f9a1").Panic(fmt.Sprintf("filesystem provider doesn't support message type %T", message))
	}

	return ioutil.WriteFile(filepath.Join(f.Path, fileName), []byte(content), 0666)
}
