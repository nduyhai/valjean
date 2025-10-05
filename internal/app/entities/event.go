package entities

type Event struct {
	ChatID            int64
	OriginalMessageId int
	ReplyMessage      string
}
