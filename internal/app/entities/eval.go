package entities

type EvalInput struct {
	ChatID       int64
	MessageID    int
	UserID       int64
	UserHandle   string
	Text         string
	ContextSnips []string // optional evidence
}

type EvalOutput struct {
	Summary    string // concise reply
	Citations  []string
	Confidence float32
}
