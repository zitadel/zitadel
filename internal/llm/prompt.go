package llm

// Prompt carries the two parts of a model conversation: the static role
// definition (System) and the per-request context data (User).
type Prompt struct {
	System string
	User   string
}
