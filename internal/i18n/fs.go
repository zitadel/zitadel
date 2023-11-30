package i18n

import (
	"github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"
	"net/http"
)

var zitadelFS, loginFS, notificationFS http.FileSystem

type Namespace string

const (
	ZITADEL      Namespace = "zitadel"
	LOGIN                  = "login"
	NOTIFICATION           = "notification"
)

func LoadFilesystem(ns Namespace) http.FileSystem {
	var err error
	switch ns {
	case ZITADEL:
		if zitadelFS != nil {
			return zitadelFS
		}
		zitadelFS, err = fs.NewWithNamespace(string(ns))
		return zitadelFS
	case LOGIN:
		if loginFS != nil {
			return loginFS
		}
		loginFS, err = fs.NewWithNamespace(string(ns))
		return loginFS
	case NOTIFICATION:
		if notificationFS != nil {
			return notificationFS
		}
		notificationFS, err = fs.NewWithNamespace(string(ns))
		return notificationFS
	}
	if err != nil {
		logging.WithFields("namespace", ns).OnError(err).Panic("unable to get namespace")
	}
	return nil
}
