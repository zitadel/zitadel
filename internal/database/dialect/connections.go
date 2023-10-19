package dialect

import "errors"

type ConnectionInfo struct {
	EventstorePusher ConnectionConfig
	ZITADEL          ConnectionConfig
}

type ConnectionConfig struct {
	MaxOpenConns,
	MaxIdleConns uint32
}

func NewConnectionInfo(openConns, idleConns uint32, pusherRatio float64) (*ConnectionInfo, error) {
	if pusherRatio < 0 || pusherRatio > 1 {
		return nil, errors.New("EventPushConnRatio must be between 0 and 1")
	}
	if openConns != 0 && openConns < 2 {
		return nil, errors.New("MaxOpenConns of the database must be higher that 1")
	}

	info := new(ConnectionInfo)

	info.EventstorePusher.MaxOpenConns = uint32(pusherRatio * float64(openConns))
	info.EventstorePusher.MaxIdleConns = uint32(pusherRatio * float64(idleConns))

	if openConns != 0 && info.EventstorePusher.MaxOpenConns < 1 && pusherRatio > 0 {
		info.EventstorePusher.MaxOpenConns = 1
	}
	if idleConns != 0 && info.EventstorePusher.MaxIdleConns < 1 && pusherRatio > 0 {
		info.EventstorePusher.MaxIdleConns = 1
	}

	if openConns != 0 {
		info.ZITADEL.MaxOpenConns = openConns - info.EventstorePusher.MaxOpenConns
	}
	if idleConns != 0 {
		info.ZITADEL.MaxIdleConns = idleConns - info.EventstorePusher.MaxIdleConns
	}

	return info, nil
}
