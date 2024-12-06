package types

type notifyResult struct {
	url                                string
	args                               map[string]interface{}
	messageType                        string
	allowUnverifiedNotificationChannel bool
}

// mockNotify returns a notifyResult and Notify function for easy mocking.
// The notifyResult will only be populated after Notify is called.
func mockNotify() (*notifyResult, Notify) {
	dst := new(notifyResult)
	return dst, func(url string, args map[string]interface{}, messageType string, allowUnverifiedNotificationChannel bool) error {
		*dst = notifyResult{
			url:                                url,
			args:                               args,
			messageType:                        messageType,
			allowUnverifiedNotificationChannel: allowUnverifiedNotificationChannel,
		}
		return nil
	}
}
