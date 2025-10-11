package entities

type Event struct {
	SourceType        SourceType
	ChatID            string
	OriginalMessageId string
	ReplyMessage      string
}
