package entities

type SourceType string

const (
	SourceTelegram SourceType = "telegram"
	SourceZalo     SourceType = "zalo"
)

type EvalInput struct {
	SourceType   SourceType
	ChatID       string
	MessageID    string
	UserID       string
	UserHandle   string
	Text         string
	ContextSnips []string // optional evidence
	ChatType     string
	ReplyFor     string
}

type EvalOutput struct {
	Summary    string // concise reply
	Citations  []string
	Confidence float32
}
