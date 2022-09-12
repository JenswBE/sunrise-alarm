package entities

import (
	"fmt"
)

const TimeFormatDate = "2006-01-02"

// Template represents a single HTML page template
type Template interface {
	GetTemplateName() string
	SetMessages(messages []Message)
}

type BaseData struct {
	Title      string
	ParentPath string
	Messages   []Message
}

// AddMessage appends a single message to the slice of messages.
// Message format and arguments are passed to fmt.Sprintf.
func (b *BaseData) AddMessage(messageType MessageType, messageFormat string, args ...any) {
	b.Messages = append(b.Messages, Message{
		Type:    messageType,
		Content: fmt.Sprintf(messageFormat, args...),
	})
}

// SetMessages replaces the current slice of messages
func (b *BaseData) SetMessages(messages []Message) {
	b.Messages = messages
}
