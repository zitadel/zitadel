package senders

// DeliveryInfo describes how an email notification was handled.
// It is part of the sent events of the notification.
type DeliveryInfo struct {
	// ProviderID is the ID of the email provider the notification was sent through.
	// It is empty if the delivery was suppressed.
	ProviderID string `json:"providerId,omitzero"`
	// DeliverySuppressed is true if the notification was accepted, but the email was not sent to the provider,
	// because the recipient domain is reserved and the operator rule of the provider suppresses them.
	DeliverySuppressed bool `json:"deliverySuppressed,omitzero"`
}
