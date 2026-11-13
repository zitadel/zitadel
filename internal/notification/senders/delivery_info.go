package senders

// DeliveryInfo describes how an email notification was handled.
// It is part of the sent events of the notification.
type DeliveryInfo struct {
	// ProviderID is the ID of the email provider the notification was sent through.
	ProviderID string `json:"providerId,omitzero"`
}

// DeliveredBy returns the filter for the payload of sent events
// of emails that were sent through the given provider.
func DeliveredBy(providerID string) map[string]any {
	return map[string]any{"providerId": providerID}
}
