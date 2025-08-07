package domain

import (
	"time"
)

type NotificationType int32

const (
	NotificationTypeEmail NotificationType = iota
	NotificationTypeSms

	notificationCount
)

type NotificationProviderState int32

const (
	NotificationProviderStateUnspecified NotificationProviderState = iota
	NotificationProviderStateActive
	NotificationProviderStateRemoved

	notificationProviderCount
)

func (s NotificationProviderState) Exists() bool {
	return s == NotificationProviderStateActive
}

type NotificationProviderType int32

const (
	NotificationProviderTypeFile NotificationProviderType = iota
	NotificationProviderTypeLog

	notificationProviderTypeCount
)

type NotificationArguments struct {
	Origin          string        `json:"origin,omitempty"`
	Domain          string        `json:"domain,omitempty"`
	Expiry          time.Duration `json:"expiry,omitempty"`
	TempUsername    string        `json:"tempUsername,omitempty"`
	ApplicationName string        `json:"applicationName,omitempty"`
	CodeID          string        `json:"codeID,omitempty"`
	SessionID       string        `json:"sessionID,omitempty"`
	AuthRequestID   string        `json:"authRequestID,omitempty"`
}

// ToMap creates a type safe map of the notification arguments.
// Since these arguments are used in text template, all keys must be PascalCase and types must remain the same (e.g. Duration).
func (n *NotificationArguments) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	if n == nil {
		return m
	}
	m["Origin"] = n.Origin
	m["Domain"] = n.Domain
	m["Expiry"] = n.Expiry
	m["TempUsername"] = n.TempUsername
	m["ApplicationName"] = n.ApplicationName
	m["CodeID"] = n.CodeID
	m["SessionID"] = n.SessionID
	m["AuthRequestID"] = n.AuthRequestID
	return m
}
