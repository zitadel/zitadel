package chat

type ChatConfig struct {
	// Defaults to true if DebugMode is set to true
	Enabled    *bool
	Url        string
	SplitCount int
	Compact    bool
}
