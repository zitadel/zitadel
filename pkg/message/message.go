package message

func (m *LocalizedMessage) LocalizationKey() string {
	return m.Key
}

func (m *LocalizedMessage) SetLocalizedMessage(message string) {
	m.Message = message
}
